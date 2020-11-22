package easymongo

type Aggregation struct{}

// TODO: Aggregation.Iter
func (p *Aggregation) Iter() *Iter { return nil }

// TODO: Aggregation.All
func (p *Aggregation) All(result interface{}) error { return nil }

// TODO: Aggregation.One
func (p *Aggregation) One(result interface{}) error { return nil }

// TODO: Aggregation.Explain
func (p *Aggregation) Explain(result interface{}) error { return nil }

// TODO: Aggregation.AllowDiskUse
func (p *Aggregation) AllowDiskUse() *Aggregation { return nil }

// TODO: Aggregation.Batch
func (p *Aggregation) Batch(n int) *Aggregation { return nil }
