package specrend

import (
	"log"
	"testing"
)

func TestAll(*testing.T) {
	cs := SMPTEsystem

	log.Printf("Temperature       x      y      z       R     G     B\n")
	log.Printf("-----------    ------ ------ ------   ----- ----- -----\n")

	for t := float64(1000); t <= 10000; t += 500 {
		xyz := SpectrumToXYZ(t, BBSpectrum)
		rgb := xyz.RGB(&cs).ConstrainRGB().NormalizeRGB()
		log.Printf("  %5.0f K      %.4f %.4f %.4f   %.3f %.3f %.3f", t, xyz.X, xyz.Y, xyz.Z, rgb.X, rgb.Y, rgb.Z)
	}
}
