package controllers
import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
    "kubequntumblock/schemas"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)


type LogStreamer struct {
	b bytes.Buffer
}

func (l *LogStreamer) String() string {
	return l.b.String()
}

/* hellogo  */
func (l *LogStreamer) Write(p []byte) (n int, err error) {
	a := strings.TrimSpace(string(p))
	l.b.WriteString(a)
	return len(p), nil
}