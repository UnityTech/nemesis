package report

// Reporter is a simple interface for publishing reports
type Reporter interface {
	Publish(reports []Report) error
}
