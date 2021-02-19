package log

import "context"

// The log is written asynchronously in this way,
// but the calling line cannot be determined.
func worker(ctx context.Context, l *Logger, ch <-chan msgOrder) {
	for {
		select {
		case m := <-ch:
			l.outPut(m.msg, m.order)
		case <-ctx.Done():
			// process all message before the ending
			for v := range ch {
				l.outPut(v.msg, v.order)
			}
			return
		default:
		}
	}
}
