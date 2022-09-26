package clients

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients/supported"
)

type ArchitectureType string

const (
	ArchitectureTypeI386 ArchitectureType = "i386"
	// TODO rename to X86_64
	ArchitectureTypeX8664      ArchitectureType = "x86_64"
	ArchitectureTypeArm64      ArchitectureType = "arm64"
	ArchitectureTypeAppleX8664 ArchitectureType = "apple-x86_64"
	ArchitectureTypeAppleArm64 ArchitectureType = "apple-arm64"
)

func (at *ArchitectureType) String() string {
	return string(*at)
}

func MapArchitectures(_ context.Context, arch string) (ArchitectureType, error) {
	switch {
	case arch == "x86_64_mac":
		return ArchitectureTypeAppleX8664, nil
	case arch == "arm64_mac":
		return ArchitectureTypeAppleArm64, nil
	case arch == "i386":
		return ArchitectureTypeI386, nil
	case arch == "x86-64" || arch == "x86_64" || arch == "x64":
		return ArchitectureTypeX8664, nil
	case arch == "aarch64" || arch == "arm64" || arch == "Arm64" || arch == "arm":
		return ArchitectureTypeArm64, nil
	}
	return "", fmt.Errorf("%s: %w", arch, supported.ErrArchitectureNotSupported)
}
