package storage

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/anacrolix/torrent/metainfo"
)

func testMarkedCompleteMissingOnRead(t *testing.T, csf func(string) ClientImplCloser) {
	td := t.TempDir()
	cic := csf(td)
	defer cic.Close()
	cs := NewClient(cic)
	info := &metainfo.Info{
		PieceLength: 1,
		Files:       []metainfo.FileInfo{{Path: []string{"a"}, Length: 1}},
	}
	ts, err := cs.OpenTorrent(info, metainfo.Hash{})
	require.NoError(t, err)
	p := ts.Piece(info.Piece(0))
	require.NoError(t, p.MarkComplete())
	// require.False(t, p.GetIsComplete())
	n, err := p.ReadAt(make([]byte, 1), 0)
	require.Error(t, err)
	require.EqualValues(t, 0, n)
	require.False(t, p.Completion().Complete)
}

func TestMarkedCompleteMissingOnReadFile(t *testing.T) {
	testMarkedCompleteMissingOnRead(t, NewFile)
}

func TestMarkedCompleteMissingOnReadFileBoltDB(t *testing.T) {
	testMarkedCompleteMissingOnRead(t, NewBoltDB)
}
