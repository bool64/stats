package stats

import "context"

type keysAndValuesCtxKey struct{}

// AddKeysAndValues returns context with added key-value pairs.
// If key-value pairs exist in parent context already, new pairs are appended.
func AddKeysAndValues(ctx context.Context, keysAndValues ...string) context.Context {
	return context.WithValue(ctx, keysAndValuesCtxKey{}, append(KeysAndValues(ctx), keysAndValues...))
}

// KeysAndValues returns key-pairs found in context or nil.
func KeysAndValues(ctx context.Context) []string {
	keysAndValues, ok := ctx.Value(keysAndValuesCtxKey{}).([]string)
	if !ok {
		return nil
	}

	return keysAndValues[:len(keysAndValues):len(keysAndValues)]
}
