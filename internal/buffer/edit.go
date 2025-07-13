package buffer

import "log/slog"

func (b *Buffer) Insert(char rune, loc Location) {
	b.cursor.X = loc.X
	b.cursor.Y = loc.Y
	b.Rows[loc.Y].Insert(loc.X, char)
	b.cursor.X++
	b.logger.Debug("insert char", slog.String("char", string(char)), slog.Any("loc", loc))
}

func (b *Buffer) Delete(loc Location) rune {
	b.cursor.X = loc.X
	b.cursor.Y = loc.Y
	char := b.Rows[loc.Y].Delete(loc.X)
	b.cursor.X--
	if b.cursor.X < 0 {
		b.cursor.X = 0
	}
	b.logger.Debug("delete char", slog.String("char", string(char)), slog.Any("loc", loc))
	return char
}
