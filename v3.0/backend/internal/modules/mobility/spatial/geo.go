package spatial

import (
	"errors"
	"math"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
)

const EarthRadiusMeters = 6371008.8

func NormalizeLongitude(lon float64) float64 {
	if lon >= -180 && lon <= 180 {
		return lon
	}
	lon = math.Mod(lon+180, 360)
	if lon < 0 {
		lon += 360
	}
	return lon - 180
}

func ValidatePosition(p model.Position) error {
	if math.IsNaN(p.Latitude) || math.IsInf(p.Latitude, 0) || p.Latitude < -90 || p.Latitude > 90 {
		return errors.New("latitude outside [-90,90]")
	}
	if math.IsNaN(p.Longitude) || math.IsInf(p.Longitude, 0) || p.Longitude < -180 || p.Longitude > 180 {
		return errors.New("longitude outside [-180,180]")
	}
	return nil
}

func DistanceMeters(a, b model.Position) float64 {
	toRad := math.Pi / 180
	lat1, lat2 := a.Latitude*toRad, b.Latitude*toRad
	dLat := (b.Latitude - a.Latitude) * toRad
	dLonDeg := NormalizeLongitude(b.Longitude - a.Longitude)
	dLon := dLonDeg * toRad
	x := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	return 2 * EarthRadiusMeters * math.Atan2(math.Sqrt(x), math.Sqrt(1-x))
}

func ValidateObservation(previous *model.SpatialState, next model.SpatialState) model.MobilityQuality {
	q := model.MobilityQuality{Valid: true, Confidence: 1}
	if err := ValidatePosition(next.ReportedPosition); err != nil {
		q.Valid = false
		q.Confidence = 0
		q.Anomalies = append(q.Anomalies, "INVALID_COORDINATES")
	}
	if next.SpeedMPS != nil && (*next.SpeedMPS < 0 || math.IsNaN(*next.SpeedMPS) || math.IsInf(*next.SpeedMPS, 0)) {
		q.Valid = false
		q.Confidence = 0
		q.Anomalies = append(q.Anomalies, "INVALID_SPEED")
	}
	if next.HeadingDegrees != nil {
		v := math.Mod(*next.HeadingDegrees, 360)
		if v < 0 {
			v += 360
		}
		next.HeadingDegrees = &v
	}
	if previous != nil && next.ObservedAt.After(previous.ObservedAt) {
		implied := DistanceMeters(previous.ReportedPosition, next.ReportedPosition) / next.ObservedAt.Sub(previous.ObservedAt).Seconds()
		limit := 90.0
		switch next.MobilityProfile {
		case model.MobilityGroundRobot:
			limit = 20
		case model.MobilityAerialDrone:
			limit = 120
		case model.MobilityStatic:
			limit = 2
		}
		if implied > limit {
			q.Confidence = .25
			q.Anomalies = append(q.Anomalies, "IMPLAUSIBLE_JUMP")
		}
		if implied < .15 {
			q.Anomalies = append(q.Anomalies, "STATIONARY_DEADBAND")
		}
	}
	return q
}
