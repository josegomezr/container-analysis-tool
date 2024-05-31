// downloader copies an image to a tarball
package downloader

import (
	"context"
	"fmt"
	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/manifest"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	"os"
	"sync"
)

type Opts struct {
	Source            string
	DestinationFolder string
}

type DownloadResult struct {
	Architecture string
	Digest       string
	Filename     string
}

func Download(ctx context.Context, opts *Opts) ([]DownloadResult, error) {
	sysctx := &types.SystemContext{}

	srcRef, err := alltransports.ParseImageName(opts.Source)
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return nil, err
	}

	src, err := srcRef.NewImageSource(ctx, sysctx)
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return nil, err
	}
	// fmt.Printf("src: %+v \n", src)

	archs, err := archsFromImageSrc(ctx, src)
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return nil, err
	}

	fmt.Printf("archs: %+v \n", archs)
	var output []DownloadResult
	var wg sync.WaitGroup

	for _, arch := range archs {
		wg.Add(1)
		go func(arch string) {
			defer wg.Done()

			result, err := DownloadArch(ctx, srcRef, opts, arch)
			if err != nil {
				fmt.Printf("err: %v \n", err)
				return
			}
			output = append(output, result)
		}(arch)
	}
	wg.Wait()

	return output, nil
}

func DownloadArch(ctx context.Context, srcRef types.ImageReference, opts *Opts, arch string) (DownloadResult, error) {
	policy, err := signature.DefaultPolicy(nil)
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return DownloadResult{}, err
	}

	sig, err := signature.NewPolicyContext(policy)
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return DownloadResult{}, err
	}

	destFile := fmt.Sprintf("%s/%s.tar", opts.DestinationFolder, arch)
	destRef, err := alltransports.ParseImageName(fmt.Sprintf("docker-archive:%s", destFile))
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return DownloadResult{}, err
	}

	err = os.RemoveAll(destFile)
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return DownloadResult{}, err
	}
	fmt.Println(transports.ImageName(destRef))

	manifestBytes, err := copy.Image(ctx, sig, destRef, srcRef, &copy.Options{
		RemoveSignatures: true,
		SourceCtx:        &types.SystemContext{ArchitectureChoice: arch},
	})

	if err != nil {
		fmt.Printf("err: %v \n", err)
		return DownloadResult{}, err
	}

	manifestDigest, err := manifest.Digest(manifestBytes)
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return DownloadResult{}, err
	}

	err = os.WriteFile(fmt.Sprintf("%s/%s.digest", opts.DestinationFolder, arch), []byte(manifestDigest.String()), 0755)
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return DownloadResult{}, err
	}

	fmt.Printf("manifestBytes: %v \n", manifestDigest.String())

	return DownloadResult{
		Architecture: arch,
		Filename:     destFile,
		Digest:       manifestDigest.String(),
	}, nil
}
