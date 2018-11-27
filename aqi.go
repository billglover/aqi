package aqi

// RGB takes an AQI reading and returns the RGB colours according to the
// scale shown on: http://aqicn.org
func RGB(aqi int) (int, int, int) {
	switch {
	case aqi <= 50:
		return 0, 152, 0
	case aqi <= 100:
		return 255, 152, 0
	case aqi <= 150:
		return 255, 50, 0
	case aqi <= 200:
		return 255, 0, 0
	case aqi <= 300:
		return 255, 0, 152
	case aqi > 300:
		return 255, 0, 35
	default:
		return 0, 0, 0
	}
}
