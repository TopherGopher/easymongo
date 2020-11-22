package easymongo

import "go.mongodb.org/mongo-driver/mongo/options"

// FindOneOptions generates the native mongo driver FindOneOptions from the FindQuery
func (q *FindQuery) FindOneOptions() *options.FindOneOptions {
	opts := options.FindOne()
	// if len(q.comment) > 0 {
	// 	opts.SetComment(q.comment)
	// }
	// if len(q.hintIndices) > 0 {
	// 	opts.SetHint(q.hintIndices)
	// }
	// if q.sortFields != 0 {
	// 	opts.SetSort(q.sortFields)
	// }
	// opts.SetAllowDiskUse(b)
	// TODO: Support the rest of the find options
	return opts
}

// FindOptions generates the native mongo driver FindOptions from the FindQuery
func (q *FindQuery) FindOptions() *options.FindOptions {
	opts := &options.FindOptions{
		AllowDiskUse:        q.allowDiskUse,
		AllowPartialResults: q.allowPartialResults,
		BatchSize:           q.batchSize,
		Collation:           q.collation,
		Comment:             q.comment,
		// CursorType:          x,
		// Hint:                x,
		Limit: q.limit,
		// Max:                 x,
		// MaxAwaitTime:        x,
		MaxTime: q.timeout,
		// Min:                 x,
		// NoCursorTimeout:     x,
		// OplogReplay:         x,
		Projection: q.projection,
		// ReturnKey:           x,
		// ShowRecordID: x,
		Skip: q.skip,
		// Snapshot:     x,
		// Sort:         x,
	}

	if len(q.hintIndices) > 0 {
		opts.SetHint(q.hintIndices)
	}
	if len(q.sortFields) != 0 {
		opts.SetSort(q.sortFields)
	}

	// TODO: Support the rest of the find options
	return opts
}

// CountOptions generates the native mongo driver CountOptions from the FindQuery
func (q *FindQuery) CountOptions() *options.CountOptions {
	opts := &options.CountOptions{
		Limit:     q.limit,
		Skip:      q.skip,
		MaxTime:   q.timeout,
		Collation: q.collation,
	}
	if len(q.hintIndices) > 0 {
		opts.SetHint(q.hintIndices)
	}
	// comments are not yet supported for count queries
	// TODO: Should we set an error on an invalid option path? Or does it matter?
	// if len(q.comment) > 0 {
	// 	opts.SetComment(q.comment)
	// }

	return opts
}
