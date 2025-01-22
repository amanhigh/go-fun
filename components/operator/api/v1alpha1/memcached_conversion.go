package v1alpha1

import (
	"fmt"

	cachev1beta1 "github.com/amanhigh/go-fun/components/operator/api/v1beta1"
	"github.com/amanhigh/go-fun/components/operator/common"

	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this Memcached to the Hub version (vbeta1).
func (src *Memcached) ConvertTo(dstRaw conversion.Hub) error {
	dst, ok := dstRaw.(*cachev1beta1.Memcached)
	if !ok {
		return fmt.Errorf("expected type *cachev1beta1.Memcached, got %T", dstRaw)
	}
	dst.Spec.Size = src.Spec.Size
	dst.Spec.ContainerPort = src.Spec.ContainerPort
	//Assume Default Sidecar Image.
	dst.Spec.SidecarImage = common.SIDECAR_IMAGE_NAME
	dst.ObjectMeta = src.ObjectMeta
	return nil
}

// ConvertFrom converts from the Hub version (vbeta1) to this version.
func (dst *Memcached) ConvertFrom(srcRaw conversion.Hub) error {
	src, ok := srcRaw.(*cachev1beta1.Memcached)
	if !ok {
		return fmt.Errorf("expected type *cachev1beta1.Memcached, got %T", srcRaw)
	}
	dst.Spec.Size = src.Spec.Size
	dst.Spec.ContainerPort = src.Spec.ContainerPort
	dst.ObjectMeta = src.ObjectMeta
	return nil
}
