package v1

import (
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"golang.org/x/net/context"

	"io"
	"io/ioutil"
	"strings"
	"testing"
)

var testConvert = Convert(func(reader io.Reader, writer io.Writer) error {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	_, err = writer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
})

var testErr = errors.New("test error")

type failingReader struct{}

func (failingReader) Read(p []byte) (n int, err error) {
	return 0, testErr
}

func TestConvert_Mutate(t *testing.T) {
	g := NewWithT(t)
	reader := strings.NewReader("test me plz")
	bytes, err := testConvert.Mutate(context.TODO(), "", reader, "")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(bytes).To(Equal([]byte("test me plz")))
}

func TestConvert_Mutate_reader_err(t *testing.T) {
	g := NewWithT(t)
	_, err := testConvert.Mutate(context.TODO(), "", failingReader{}, "")
	g.Expect(err).To(HaveOccurred())
}
