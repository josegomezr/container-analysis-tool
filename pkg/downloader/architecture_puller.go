package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/containers/image/v5/types"
)

func archsFromBytes(instream []byte) []string {
	var manistruct struct {
		Manifests []struct {
			Platform struct {
				Architecture string `json:"architecture"`
			} `json:"platform"`
		} `json:"manifests"`
	}

	err := json.Unmarshal(instream, &manistruct)

	if err != nil {
		return []string{}
	}

	var archs []string
	for _, descriptor := range manistruct.Manifests {
		archs = append(archs, descriptor.Platform.Architecture)
	}

	return archs
}

func archsFromImageSrc(ctx context.Context, src types.ImageSource) ([]string, error) {
	man, mediaType, err := src.GetManifest(ctx, nil)
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return nil, err
	}
	// fmt.Printf("man: %v \n", man)
	fmt.Printf("mediaType: %v \n", mediaType)
	return archsFromBytes(man), nil
}
