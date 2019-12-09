package assetstore

import (
	"fmt"
	"github.com/kyma-project/kyma/components/application-registry/internal/apperrors"
	"github.com/kyma-project/kyma/components/application-registry/internal/metadata/specification/assetstore/docstopic"
	"github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/retry"
)

const (
	DocsTopicModeSingle = "single"
	DocsTopicNameFormat = "%s-%s"
)

type ResourceInterface interface {
	Get(name string, opts metav1.GetOptions, subresources ...string) (*unstructured.Unstructured, error)
	Delete(name string, opts *metav1.DeleteOptions, subresources ...string) error
	Create(obj *unstructured.Unstructured, options metav1.CreateOptions, subresources ...string) (*unstructured.Unstructured, error)
	Update(obj *unstructured.Unstructured, options metav1.UpdateOptions, subresources ...string) (*unstructured.Unstructured, error)
}

type DocsTopicRepository interface {
	Get(id string) (docstopic.Entry, apperrors.AppError)
	Upsert(documentationTopic docstopic.Entry) apperrors.AppError
	Delete(id string) apperrors.AppError
}

type repository struct {
	resourceInterface ResourceInterface
}

func NewDocsTopicRepository(resourceInterface ResourceInterface) DocsTopicRepository {
	return repository{
		resourceInterface: resourceInterface,
	}
}

func (r repository) Upsert(docsTopicEntry docstopic.Entry) apperrors.AppError {
	_, err := r.get(docsTopicEntry.Id)
	if err != nil && err.Code() == apperrors.CodeNotFound {
		return r.create(toK8sType(docsTopicEntry))
	}

	if err != nil {
		return err
	}

	k8sDocsTopic := toK8sType(docsTopicEntry)

	return r.update(docsTopicEntry.Id, k8sDocsTopic)
}

func (r repository) Get(id string) (docstopic.Entry, apperrors.AppError) {
	docsTopic, err := r.get(id)
	if err != nil {
		return docstopic.Entry{}, err
	}

	return fromK8sType(docsTopic), nil
}

func (r repository) Delete(id string) apperrors.AppError {
	err := r.resourceInterface.Delete(id, &metav1.DeleteOptions{})
	if err != nil && !k8serrors.IsNotFound(err) {
		return apperrors.Internal("Failed to delete DocsTopic: %s.", err)
	}

	return nil
}

func (r repository) get(id string) (v1beta1.ClusterAssetGroup, apperrors.AppError) {
	u, err := r.resourceInterface.Get(id, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return v1beta1.ClusterAssetGroup{}, apperrors.NotFound("Docs Topic with %s id not found.", id)
		}

		return v1beta1.ClusterAssetGroup{}, apperrors.Internal("Failed to get Docs Topic, %s.", err)
	}

	var docsTopic v1beta1.ClusterAssetGroup
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &docsTopic)
	if err != nil {
		return v1beta1.ClusterAssetGroup{}, apperrors.Internal("Failed to convert from unstructured object, %s.", err)
	}

	return docsTopic, nil
}

func (r repository) create(docsTopic v1beta1.ClusterAssetGroup) apperrors.AppError {
	u, err := toUstructured(docsTopic)
	if err != nil {
		return apperrors.Internal("Failed to create Documentation Topic, %s.", err)
	}

	_, err = r.resourceInterface.Create(u, metav1.CreateOptions{})
	if err != nil {
		return apperrors.Internal("Failed to create Documentation Topic, %s.", err)
	}

	return nil
}

func (r repository) update(id string, docsTopic v1beta1.ClusterAssetGroup) apperrors.AppError {

	getRefreshedDocsTopic := func(id string, docsTopic v1beta1.ClusterAssetGroup) (v1beta1.ClusterAssetGroup, error) {
		newUnstructured, err := r.resourceInterface.Get(id, metav1.GetOptions{})
		if err != nil {
			return v1beta1.ClusterAssetGroup{}, err
		}

		newDocsTopic, err := fromUnstructured(newUnstructured)
		if err != nil {
			return v1beta1.ClusterAssetGroup{}, err
		}

		newDocsTopic.Spec = docsTopic.Spec

		return newDocsTopic, nil
	}

	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		newDocsTopic, err := getRefreshedDocsTopic(id, docsTopic)
		if err != nil {
			return err
		}

		u, err := toUstructured(newDocsTopic)
		if err != nil {
			return err
		}

		_, err = r.resourceInterface.Update(u, metav1.UpdateOptions{})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return apperrors.Internal("Failed to update Documentation Topic, %s.", err)
	}

	return nil
}

func toUstructured(docsTopic v1beta1.ClusterAssetGroup) (*unstructured.Unstructured, error) {
	obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&docsTopic)
	if err != nil {
		return nil, err
	}

	return &unstructured.Unstructured{Object: obj}, nil
}

func fromUnstructured(u *unstructured.Unstructured) (v1beta1.ClusterAssetGroup, error) {
	var docsTopic v1beta1.ClusterAssetGroup
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &docsTopic)
	if err != nil {
		return v1beta1.ClusterAssetGroup{}, err
	}

	return docsTopic, nil
}

func toK8sType(docsTopicEntry docstopic.Entry) v1beta1.ClusterAssetGroup {
	sources := make([]v1beta1.Source, 0, 3)
	for key, url := range docsTopicEntry.Urls {
		source := v1beta1.Source{
			Name: v1beta1.AssetGroupSourceName(fmt.Sprintf(DocsTopicNameFormat, key, docsTopicEntry.Id)),
			URL:  url,
			Mode: DocsTopicModeSingle,
			Type: v1beta1.AssetGroupSourceType(key),
		}
		sources = append(sources, source)
	}

	return v1beta1.ClusterAssetGroup{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterDocsTopic",
			APIVersion: v1beta1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      docsTopicEntry.Id,
			Namespace: "kyma-integration",
			Labels:    docsTopicEntry.Labels,
		},
		Spec: v1beta1.ClusterAssetGroupSpec{
			CommonAssetGroupSpec: v1beta1.CommonAssetGroupSpec{
				DisplayName: "Some display name",
				Description: "Some description",
				Sources:     sources,
			},
		}}
}

func fromK8sType(k8sDocsTopic v1beta1.ClusterAssetGroup) docstopic.Entry {
	urls := make(map[string]string)

	for _, source := range k8sDocsTopic.Spec.Sources {
		urls[string(source.Type)] = source.URL
	}

	return docstopic.Entry{
		Id:          k8sDocsTopic.Name,
		Description: k8sDocsTopic.Spec.Description,
		DisplayName: k8sDocsTopic.Spec.DisplayName,
		Urls:        urls,
		Labels:      k8sDocsTopic.Labels,
		Status:      docstopic.StatusType(k8sDocsTopic.Status.Phase),
	}
}
