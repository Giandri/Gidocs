package jobs

import "testing"

func TestRunUnknownTool(t *testing.T) {
	err := Run(t.Context(), Payload{JobID: NewID(), Tool: "merge-pdf", OutDir: t.TempDir()})
	if err == nil {
		t.Fatal("Run dengan tool tak dikenal harus mengembalikan error")
	}
}
