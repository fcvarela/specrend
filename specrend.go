// Package specrend provides an implementation of John Walker's 'Color Rendering of Spectra' methods.
// See original documentation and reference implementation at: http://www.fourmilab.ch/ and http://www.fourmilab.ch/documents/specrend/
package specrend

import "math"

// Vec2d is a two-dimensional float64 vector
type Vec2d struct {
	X float64
	Y float64
}

// Vec3d is a three-dimensional float64 vector
type Vec3d struct {
	X float64
	Y float64
	Z float64
}

// ColorSystem is defined by the CIE x and y coordinates of
// its three primary illuminants and the x and y coordinates of
// the white point.
type ColorSystem struct {
	Name  string
	Red   Vec2d
	Green Vec2d
	Blue  Vec2d
	White Vec2d
	Gamma float64
}

// White point chromaticities
var (
	IlluminantC   = Vec2d{0.3101, 0.3162}
	IlluminantD65 = Vec2d{0.3127, 0.3291}
	IlluminantE   = Vec2d{1.0 / 3.0, 1.0 / 3.0}
	GammaRec709   = float64(0.0)
)

// Colorsystems
var (
	NTSCSystem   = ColorSystem{"NTSC", Vec2d{0.67, 0.33}, Vec2d{0.21, 0.71}, Vec2d{0.14, 0.08}, IlluminantC, GammaRec709}
	EBUSystem    = ColorSystem{"EBU", Vec2d{0.64, 0.33}, Vec2d{0.29, 0.60}, Vec2d{0.15, 0.06}, IlluminantD65, GammaRec709}
	SMPTEsystem  = ColorSystem{"SMPTE", Vec2d{0.630, 0.340}, Vec2d{0.310, 0.595}, Vec2d{0.155, 0.070}, IlluminantD65, GammaRec709}
	HDTVsystem   = ColorSystem{"HDTV", Vec2d{0.670, 0.330}, Vec2d{0.210, 0.710}, Vec2d{0.150, 0.060}, IlluminantD65, GammaRec709}
	CIEsystem    = ColorSystem{"CIE", Vec2d{0.7355, 0.2645}, Vec2d{0.2658, 0.7243}, Vec2d{0.1669, 0.0085}, IlluminantE, GammaRec709}
	Rec709system = ColorSystem{"CIE REC 709", Vec2d{0.64, 0.33}, Vec2d{0.30, 0.60}, Vec2d{0.15, 0.06}, IlluminantD65, GammaRec709}
)

// spectum to xyz aux table
var cieColourMatch = [81][3]float64{
	{0.0014, 0.0000, 0.0065}, {0.0022, 0.0001, 0.0105}, {0.0042, 0.0001, 0.0201},
	{0.0076, 0.0002, 0.0362}, {0.0143, 0.0004, 0.0679}, {0.0232, 0.0006, 0.1102},
	{0.0435, 0.0012, 0.2074}, {0.0776, 0.0022, 0.3713}, {0.1344, 0.0040, 0.6456},
	{0.2148, 0.0073, 1.0391}, {0.2839, 0.0116, 1.3856}, {0.3285, 0.0168, 1.6230},
	{0.3483, 0.0230, 1.7471}, {0.3481, 0.0298, 1.7826}, {0.3362, 0.0380, 1.7721},
	{0.3187, 0.0480, 1.7441}, {0.2908, 0.0600, 1.6692}, {0.2511, 0.0739, 1.5281},
	{0.1954, 0.0910, 1.2876}, {0.1421, 0.1126, 1.0419}, {0.0956, 0.1390, 0.8130},
	{0.0580, 0.1693, 0.6162}, {0.0320, 0.2080, 0.4652}, {0.0147, 0.2586, 0.3533},
	{0.0049, 0.3230, 0.2720}, {0.0024, 0.4073, 0.2123}, {0.0093, 0.5030, 0.1582},
	{0.0291, 0.6082, 0.1117}, {0.0633, 0.7100, 0.0782}, {0.1096, 0.7932, 0.0573},
	{0.1655, 0.8620, 0.0422}, {0.2257, 0.9149, 0.0298}, {0.2904, 0.9540, 0.0203},
	{0.3597, 0.9803, 0.0134}, {0.4334, 0.9950, 0.0087}, {0.5121, 1.0000, 0.0057},
	{0.5945, 0.9950, 0.0039}, {0.6784, 0.9786, 0.0027}, {0.7621, 0.9520, 0.0021},
	{0.8425, 0.9154, 0.0018}, {0.9163, 0.8700, 0.0017}, {0.9786, 0.8163, 0.0014},
	{1.0263, 0.7570, 0.0011}, {1.0567, 0.6949, 0.0010}, {1.0622, 0.6310, 0.0008},
	{1.0456, 0.5668, 0.0006}, {1.0026, 0.5030, 0.0003}, {0.9384, 0.4412, 0.0002},
	{0.8544, 0.3810, 0.0002}, {0.7514, 0.3210, 0.0001}, {0.6424, 0.2650, 0.0000},
	{0.5419, 0.2170, 0.0000}, {0.4479, 0.1750, 0.0000}, {0.3608, 0.1382, 0.0000},
	{0.2835, 0.1070, 0.0000}, {0.2187, 0.0816, 0.0000}, {0.1649, 0.0610, 0.0000},
	{0.1212, 0.0446, 0.0000}, {0.0874, 0.0320, 0.0000}, {0.0636, 0.0232, 0.0000},
	{0.0468, 0.0170, 0.0000}, {0.0329, 0.0119, 0.0000}, {0.0227, 0.0082, 0.0000},
	{0.0158, 0.0057, 0.0000}, {0.0114, 0.0041, 0.0000}, {0.0081, 0.0029, 0.0000},
	{0.0058, 0.0021, 0.0000}, {0.0041, 0.0015, 0.0000}, {0.0029, 0.0010, 0.0000},
	{0.0020, 0.0007, 0.0000}, {0.0014, 0.0005, 0.0000}, {0.0010, 0.0004, 0.0000},
	{0.0007, 0.0002, 0.0000}, {0.0005, 0.0002, 0.0000}, {0.0003, 0.0001, 0.0000},
	{0.0002, 0.0001, 0.0000}, {0.0002, 0.0001, 0.0000}, {0.0001, 0.0000, 0.0000},
	{0.0001, 0.0000, 0.0000}, {0.0001, 0.0000, 0.0000}, {0.0000, 0.0000, 0.0000},
}

// XY will determine 1931 chromaticities x, y from 1976 coordinates u', v'
func (src Vec2d) XY() Vec2d {
	return Vec2d{
		(9 * src.X) / ((6 * src.X) - (16 * src.Y) + 12),
		(9 * src.Y) / ((6 * src.Y) - (16 * src.Y) + 12),
	}
}

// UpVp will determine 1976 coordinates u', v' from 1931 chromaticities x, y
func (src Vec2d) UpVp() Vec2d {
	return Vec2d{
		(4 * src.X) / ((-2 * src.X) + (12 * src.Y) + 3),
		(9 * src.Y) / ((-2 * src.X) + (12 * src.Y) + 3),
	}

}

// RGB computes RGB components from XYZ on a ColorSystem:
// Given an additive tricolour system CS, defined by the CIE x
// and y chromaticities of its three primaries (z is derived
// trivially as 1-(x+y)), and a desired chromaticity (XC, YC,
// ZC) in CIE space, determine the contribution of each
// primary in a linear combination which sums to the desired
// chromaticity.  If the  requested chromaticity falls outside
// the Maxwell  triangle (colour gamut) formed by the three
// primaries, one of the r, g, or b weights will be negative.
//
// Caller can use Vec3d.ConstrainRGB() to desaturate an
// outside-gamut colour to the closest representation within
// the available gamut and/or norm_rgb to normalise the RGB
// components so the largest nonzero component has value 1.
func (src Vec3d) RGB(cs *ColorSystem) Vec3d {
	var xr, yr, zr, xg, yg, zg, xb, yb, zb float64
	var xw, yw, zw float64
	var rx, ry, rz, gx, gy, gz, bx, by, bz float64
	var rw, gw, bw float64

	xr = cs.Red.X
	yr = cs.Red.Y
	zr = 1.0 - (xr + yr)
	xg = cs.Green.X
	yg = cs.Green.Y
	zg = 1.0 - (xg + yg)
	xb = cs.Blue.X
	yb = cs.Blue.Y
	zb = 1.0 - (xb + yb)

	xw = cs.White.X
	yw = cs.White.Y
	zw = 1 - (xw + yw)

	// xyz . rgb matrix, before scaling to white.
	rx = (yg * zb) - (yb * zg)
	ry = (xb * zg) - (xg * zb)
	rz = (xg * yb) - (xb * yg)
	gx = (yb * zr) - (yr * zb)
	gy = (xr * zb) - (xb * zr)
	gz = (xb * yr) - (xr * yb)
	bx = (yr * zg) - (yg * zr)
	by = (xg * zr) - (xr * zg)
	bz = (xr * yg) - (xg * yr)

	// White scaling factors. Dividing by yw scales the white luminance to unity, as conventional.
	rw = ((rx * xw) + (ry * yw) + (rz * zw)) / yw
	gw = ((gx * xw) + (gy * yw) + (gz * zw)) / yw
	bw = ((bx * xw) + (by * yw) + (bz * zw)) / yw

	// xyz . rgb matrix, correctly scaled to white.
	rx = rx / rw
	ry = ry / rw
	rz = rz / rw
	gx = gx / gw
	gy = gy / gw
	gz = gz / gw
	bx = bx / bw
	by = by / bw
	bz = bz / bw

	/* rgb of the desired point */
	return Vec3d{
		(rx * src.X) + (ry * src.Y) + (rz * src.Z),
		(gx * src.X) + (gy * src.Y) + (gz * src.Z),
		(bx * src.X) + (by * src.Y) + (bz * src.Z),
	}
}

// InsideGamut tests whether a requested colour is within the gamut
// achievable with the primaries of the current colour
// system.  This amounts simply to testing whether all the
// primary weights are non-negative.
func (src Vec3d) InsideGamut() bool {
	if src.X >= 0 && src.Y >= 0 && src.Z >= 0.0 {
		return true
	}
	return false
}

// ConstrainRGB desaturates a color until it lies inside a color gamut
// If the requested RGB shade contains a negative weight for
// one of the primaries, it lies outside the colour gamut
// accessible from the given triple of primaries.  Desaturate
// it by adding white, equal quantities of R, G, and B, enough
// to make RGB all positive.
func (src Vec3d) ConstrainRGB() Vec3d {
	var w float64

	/* Amount of white needed is w = - min(0, *r, *g, *b) */
	if 0.0 < src.X {
		w = 0.0
	} else {
		w = src.X
	}

	if w >= src.Y {
		w = src.Y
	}

	if w >= src.Z {
		w = src.Z
	}

	w = -w

	/* Add just enough white to make r, g, b all positive. */
	if w > 0 {
		return Vec3d{src.X + w, src.Y + w, src.Z + w}
	}
	return src
}

// Corrects a single color component
func gammaCorrectColorComponent(cs *ColorSystem, c float64) float64 {
	gamma := cs.Gamma

	if gamma == GammaRec709 {
		// Rec. 709 gamma correction.
		cc := 0.018
		if c < cc {
			return ((1.099 * math.Pow(cc, 0.45)) - 0.099) / cc
		}
		return (1.099 * math.Pow(c, 0.45)) - 0.099
	}
	// Nonlinear colour = (Linear colour)^(1/gamma)
	return math.Pow(c, 1.0/gamma)
}

// GammaCorrect transforms linear RGB values to nonlinear
func (src Vec3d) GammaCorrect(cs *ColorSystem) Vec3d {
	return Vec3d{
		gammaCorrectColorComponent(cs, src.X),
		gammaCorrectColorComponent(cs, src.Y),
		gammaCorrectColorComponent(cs, src.Z),
	}

}

// NormalizeRGB normalises RGB components so the most intense (unless all are zero) has a value of 1.
func (src Vec3d) NormalizeRGB() Vec3d {
	greatest := math.Max(src.X, math.Max(src.Y, src.Z))

	if greatest > 0.0 {
		return Vec3d{src.X / greatest, src.Y / greatest, src.Z / greatest}
	}

	return src
}

// SpectrumToXYZ computes XYZ CIE coordinates corresponding to
// a light source with spectral distribution given by  the
// function SPEC_INTENS, which is called with a series of
// wavelengths between 380 and 780 nm (the argument is
// expressed in meters), which returns emittance at  that
// wavelength in arbitrary units.  The chromaticity
// coordinates of the spectrum are returned in the x, y, and z
// arguments which respect the identity: x + y + z = 1.
func SpectrumToXYZ(temperature float64, specIntens func(temperature float64, wavelength float64) float64) Vec3d {
	var x, y, z float64

	for i, lambda := 0, 380.0; lambda < 780.1; i, lambda = i+1, lambda+5 {
		Me := specIntens(temperature, lambda)
		x += Me * cieColourMatch[i][0]
		y += Me * cieColourMatch[i][1]
		z += Me * cieColourMatch[i][2]
	}
	xyz := (x + y + z)
	return Vec3d{
		x / xyz,
		y / xyz,
		z / xyz,
	}
}

// BBSpectrum calculates, by Planck's radiation law, the emittance
// of a black body of temperature bbTemp at the given wavelength (in metres).
func BBSpectrum(temperature float64, wavelength float64) float64 {
	wlm := wavelength * 1e-9
	return (3.74183e-16 * math.Pow(wlm, -5.0)) / (math.Exp(1.4388e-2/(wlm*temperature)) - 1.0)
}
