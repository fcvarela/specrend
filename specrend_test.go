package specrend

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedXYZ map[float64]Vec3d
var expectedRGB map[float64]Vec3d

// prepare fixtures for testing against original implementation results
func init() {
	expectedXYZ = make(map[float64]Vec3d)
	expectedRGB = make(map[float64]Vec3d)

	expectedXYZ[1000] = Vec3d{0.6528, 0.3444, 0.0028}
	expectedXYZ[1500] = Vec3d{0.5857, 0.3931, 0.0212}
	expectedXYZ[2000] = Vec3d{0.5267, 0.4133, 0.0600}
	expectedXYZ[2500] = Vec3d{0.4770, 0.4137, 0.1093}
	expectedXYZ[3000] = Vec3d{0.4369, 0.4041, 0.1590}
	expectedXYZ[3500] = Vec3d{0.4053, 0.3907, 0.2040}
	expectedXYZ[4000] = Vec3d{0.3805, 0.3768, 0.2428}
	expectedXYZ[4500] = Vec3d{0.3608, 0.3636, 0.2756}
	expectedXYZ[5000] = Vec3d{0.3451, 0.3516, 0.3032}
	expectedXYZ[5500] = Vec3d{0.3325, 0.3411, 0.3265}
	expectedXYZ[6000] = Vec3d{0.3221, 0.3318, 0.3461}
	expectedXYZ[6500] = Vec3d{0.3135, 0.3237, 0.3628}
	expectedXYZ[7000] = Vec3d{0.3064, 0.3166, 0.3770}
	expectedXYZ[7500] = Vec3d{0.3004, 0.3103, 0.3893}
	expectedXYZ[8000] = Vec3d{0.2952, 0.3048, 0.4000}
	expectedXYZ[8500] = Vec3d{0.2908, 0.3000, 0.4093}
	expectedXYZ[9000] = Vec3d{0.2869, 0.2956, 0.4174}
	expectedXYZ[9500] = Vec3d{0.2836, 0.2918, 0.4246}
	expectedXYZ[10000] = Vec3d{0.2807, 0.2884, 0.4310}

	expectedRGB[1000] = Vec3d{1.000, 0.007, 0.000}
	expectedRGB[1500] = Vec3d{1.000, 0.126, 0.000}
	expectedRGB[2000] = Vec3d{1.000, 0.234, 0.010}
	expectedRGB[2500] = Vec3d{1.000, 0.349, 0.067}
	expectedRGB[3000] = Vec3d{1.000, 0.454, 0.151}
	expectedRGB[3500] = Vec3d{1.000, 0.549, 0.254}
	expectedRGB[4000] = Vec3d{1.000, 0.635, 0.370}
	expectedRGB[4500] = Vec3d{1.000, 0.710, 0.493}
	expectedRGB[5000] = Vec3d{1.000, 0.778, 0.620}
	expectedRGB[5500] = Vec3d{1.000, 0.837, 0.746}
	expectedRGB[6000] = Vec3d{1.000, 0.890, 0.869}
	expectedRGB[6500] = Vec3d{1.000, 0.937, 0.988}
	expectedRGB[7000] = Vec3d{0.907, 0.888, 1.000}
	expectedRGB[7500] = Vec3d{0.827, 0.839, 1.000}
	expectedRGB[8000] = Vec3d{0.762, 0.800, 1.000}
	expectedRGB[8500] = Vec3d{0.711, 0.766, 1.000}
	expectedRGB[9000] = Vec3d{0.668, 0.738, 1.000}
	expectedRGB[9500] = Vec3d{0.632, 0.714, 1.000}
	expectedRGB[10000] = Vec3d{0.602, 0.693, 1.000}
}

func TestRGB(t *testing.T) {
	cs := SMPTEsystem

	for temp, refRGB := range expectedRGB {
		xyz := SpectrumToXYZ(temp, BBSpectrum)
		rgb := xyz.RGB(&cs).ConstrainRGB().NormalizeRGB()

		tstVal := fmt.Sprintf("%.3f %.3f %.3f", rgb.X, rgb.Y, rgb.Z)
		refVal := fmt.Sprintf("%.3f %.3f %.3f", refRGB.X, refRGB.Y, refRGB.Z)
		assert.Equal(t, tstVal, refVal, "")
	}
}

func TestXYZ(t *testing.T) {
	for temp, refXYZ := range expectedXYZ {
		xyz := SpectrumToXYZ(temp, BBSpectrum)
		diffX, diffY, diffZ := math.Abs(refXYZ.X-xyz.X), math.Abs(refXYZ.Y-xyz.Y), math.Abs(refXYZ.Z-xyz.Z)
		if diffX > 0.0001 || diffY > 0.0001 || diffZ > 0.0001 {
			assert.Fail(t, "Computed XYZ does not match expected value")
		}
	}
}
