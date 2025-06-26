package logger

// func WithLogRequestPath(ctx context.Context, path string) context.Context {
// 	if c, ok := ctx.Value(key).(logCtx); ok {
// 		c.reqPath = path
// 		return context.WithValue(ctx, key, c)
// 	}
// 	return context.WithValue(ctx, key, logCtx{reqPath: path})
// }

// func WithLogRequestMethod(ctx context.Context, method string) context.Context {
// 	if c, ok := ctx.Value(key).(logCtx); ok {
// 		c.reqMethod = method
// 		return context.WithValue(ctx, key, c)
// 	}
// 	return context.WithValue(ctx, key, logCtx{reqMethod: method})
// }
