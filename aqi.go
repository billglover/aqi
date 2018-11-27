package aqi

// RGB takes an AQI reading and returns the RGB colours according to the
// scale shown on: http://aqicn.org
func RGB(aqi int) (int, int, int) {
	switch {
	case aqi <= 50:
		return 0, 153, 102
	case aqi <= 100:
		return 255, 222, 51
	case aqi <= 150:
		return 255, 153, 51
	case aqi <= 200:
		return 204, 0, 51
	case aqi <= 300:
		return 102, 0, 153
	case aqi > 300:
		return 126, 0, 35
	default:
		return 0, 0, 0
	}
}
