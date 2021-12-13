package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	label_user      = "user"
	label_accepted  = "Accepted"
	lablel_rejected = "Rejected"
	label_challenge = "Chalenge"
)

var (
	userStatus_metric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "radius_response",
		Help: "status of radius repsonse include (Accepted, Rejected, Chalenge)",
	}, []string{"user", "state"})
)

func make_metrics() {
	userStatus_metric.Reset()

	success := getMertricMap(&Accepted_users)
	for user, count := range success {
		userStatus_metric.WithLabelValues(user, label_accepted).Set(float64(count))
	}

	reject := getMertricMap(&Rejected_users)
	for user, count := range reject {
		userStatus_metric.WithLabelValues(user, lablel_rejected).Set(float64(count))
	}

	challenge := getMertricMap(&Chalenged_users)
	for user, count := range challenge {
		userStatus_metric.WithLabelValues(user, label_challenge).Set(float64(count))
	}

}

func getMertricMap(storage users_storage) map[string]int {

	metricMap := convertSliceToMap(storage)
	return metricMap

}

func convertSliceToMap(storage users_storage) map[string]int {
	usersActivity := make(map[string]int)
	for _, user := range storage.ReadAndDelete() {

		// count users attempts
		usersActivity[user]++
	}
	return usersActivity
}
