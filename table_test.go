package table

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// prefixedRow renders the selected row with a `> ` prefix and no styles.
type prefixedRow []interface{}

func (row prefixedRow) Render(w io.Writer, model Model, index int) {
	cells := make([]string, len(row))
	for i, v := range row {
		cells[i] = fmt.Sprintf("%v", v)
	}
	s := strings.Join(cells, "\t")
	if index == model.Cursor() {
		s = "> " + s
	} else {
		s = "  " + s
	}
	fmt.Fprintln(w, s)
}

var uncoloredStyles = Styles{
	Title:       lipgloss.NewStyle(),
	SelectedRow: lipgloss.NewStyle(),
}

func TestModel_View(t *testing.T) {
	model := New([]string{"  ID", "EMAIL", "USERNAME", "CREATED-AT"}, 0, 5)
	model.Styles = uncoloredStyles
	model.SetRows([]Row{
		prefixedRow{"1", "john@example.org", "john", "2022-02-21T18:02:29.762Z"},
		prefixedRow{"2", "bob@example.org", "bob", "2022-02-21T18:02:29.762Z"},
		prefixedRow{"3", "alice@example.org", "alice", "2022-02-21T18:02:29.762Z"},
		prefixedRow{"4", "thomas@example.org", "thomas", "2022-02-21T18:02:29.762Z"},
	})
	got := model.View()
	wantEq(t, ""+
		"  ID EMAIL              USERNAME CREATED-AT              \n"+
		"> 1  john@example.org   john     2022-02-21T18:02:29.762Z\n"+
		"  2  bob@example.org    bob      2022-02-21T18:02:29.762Z\n"+
		"  3  alice@example.org  alice    2022-02-21T18:02:29.762Z\n"+
		"  4  thomas@example.org thomas   2022-02-21T18:02:29.762Z", got)
}

func TestModel_Movements(t *testing.T) {
	model := New([]string{"  #"}, 0, 4)
	model.Styles = uncoloredStyles
	rows := make([]Row, 10)
	for i := 0; i < 10; i++ {
		rows[i] = prefixedRow{i}
	}
	model.SetRows(rows)

	initial := model.View()
	wantEq(t, ""+
		"  #\n"+
		"> 0\n"+
		"  1\n"+
		"  2", initial)

	t.Run("up", func(t *testing.T) {
		model.GoTop()
		model.GoUp()

		upFromTop := model.View()
		wantEq(t, initial, upFromTop)

		model.GoBottom()
		model.GoUp()

		upFromBottom := model.View()
		wantEq(t, ""+
			"  #\n"+
			"  7\n"+
			"> 8\n"+
			"  9", upFromBottom)

		model.GoTop()
		model.GoPageDown()
		model.GoUp()

		upFromPageDown := model.View()
		wantEq(t, ""+
			"  #\n"+
			"> 2\n"+
			"  3\n"+
			"  4", upFromPageDown)
	})

	t.Run("down", func(t *testing.T) {
		model.GoTop()
		model.GoDown()

		downFromTop := model.View()
		wantEq(t, ""+
			"  #\n"+
			"  0\n"+
			"> 1\n"+
			"  2", downFromTop)

		model.GoBottom()
		model.GoDown()

		downFromBottom := model.View()
		wantEq(t, ""+
			"  #\n"+
			"  7\n"+
			"  8\n"+
			"> 9", downFromBottom)

		model.GoBottom()
		model.GoPageUp()
		model.GoDown()

		downFromPageUp := model.View()
		wantEq(t, ""+
			"  #\n"+
			"  5\n"+
			"  6\n"+
			"> 7", downFromPageUp)

		model.GoBottom()
		model.GoPageUp()
		model.GoDown()
	})

	t.Run("pgdown", func(t *testing.T) {
		model.GoTop()
		model.GoPageDown()

		pageDownFromTop := model.View()
		wantEq(t, ""+
			"  #\n"+
			"> 3\n"+
			"  4\n"+
			"  5", pageDownFromTop)

		model.GoBottom()
		model.GoPageDown()

		pageDownFromBottom := model.View()
		wantEq(t, ""+
			"  #\n"+
			"  7\n"+
			"  8\n"+
			"> 9", pageDownFromBottom)

		model.GoUp()
		model.GoPageDown()

		pageDownFromUp := model.View()
		wantEq(t, pageDownFromBottom, pageDownFromUp)
	})

	t.Run("pgup", func(t *testing.T) {
		model.GoTop()
		model.GoPageUp()

		pageUpFromTop := model.View()
		wantEq(t, ""+
			"  #\n"+
			"> 0\n"+
			"  1\n"+
			"  2", pageUpFromTop)

		model.GoDown()
		model.GoPageUp()

		pageUpFromDown := model.View()
		wantEq(t, pageUpFromTop, pageUpFromDown)

		model.GoBottom()
		model.GoPageUp()

		pageUpFromBottom := model.View()
		wantEq(t, ""+
			"  #\n"+
			"  4\n"+
			"  5\n"+
			"> 6", pageUpFromBottom)
	})

	t.Run("home", func(t *testing.T) {
		model.GoTop()

		top := model.View()
		wantEq(t, ""+
			"  #\n"+
			"> 0\n"+
			"  1\n"+
			"  2", top)

		model.GoTop()

		topFromTop := model.View()
		wantEq(t, top, topFromTop)
	})

	t.Run("end", func(t *testing.T) {
		model.GoBottom()

		bottom := model.View()
		wantEq(t, ""+
			"  #\n"+
			"  7\n"+
			"  8\n"+
			"> 9", bottom)

		model.GoBottom()

		bottomFromBottom := model.View()
		wantEq(t, bottom, bottomFromBottom)
	})
}

func TestModel_SetSize(t *testing.T) {
	model := New([]string{"  #"}, 0, 4)
	model.Styles = uncoloredStyles
	rows := make([]Row, 10)
	for i := 0; i < 10; i++ {
		rows[i] = prefixedRow{fmt.Sprintf("item %d", i)}
	}
	model.SetRows(rows)

	initial := model.View()
	wantEq(t, ""+
		"  #     \n"+
		"> item 0\n"+
		"  item 1\n"+
		"  item 2", initial)

	model.SetSize(4, 5)

	got := model.View()
	wantEq(t, ""+
		"  # \n"+
		"> it\n"+
		"  it\n"+
		"  it\n"+
		"  it", got)

	model.GoBottom()
	model.SetSize(0, 3)

	got = model.View()

	// TODO: maybe change behavoir and keep scroll position on the selected item.
	// Instead of moving selection to the bound of the new size.
	wantEq(t, ""+
		"  #     \n"+
		"  item 6\n"+
		"> item 7", got)
}

func TestModel_SelectedRow(t *testing.T) {
	model := New([]string{"  #"}, 0, 4)
	rows := make([]Row, 10)
	for i := 0; i < 10; i++ {
		rows[i] = SimpleRow{i}
	}
	model.SetRows(rows)

	got := model.SelectedRow()
	wantEq(t, rows[0], got)

	model.GoPageDown()
	got = model.SelectedRow()
	wantEq(t, rows[3], got)
}

func TestModel_Update(t *testing.T) {
	tt := []struct {
		msg        tea.KeyMsg
		wantCursor int
	}{
		{
			msg:        tea.KeyMsg{Type: tea.KeyEnd},
			wantCursor: 9,
		},
		{
			msg:        tea.KeyMsg{Type: tea.KeyHome},
			wantCursor: 0,
		},
		{
			msg:        tea.KeyMsg{Type: tea.KeyPgDown},
			wantCursor: 3,
		},
		{
			msg:        tea.KeyMsg{Type: tea.KeyPgUp},
			wantCursor: 0,
		},
		{
			msg:        tea.KeyMsg{Type: tea.KeyDown},
			wantCursor: 1,
		},
		{
			msg:        tea.KeyMsg{Type: tea.KeyUp},
			wantCursor: 0,
		},
	}
	for _, tc := range tt {
		t.Run(tc.msg.String(), func(t *testing.T) {
			model := New([]string{"  #"}, 0, 4)
			rows := make([]Row, 10)
			for i := 0; i < 10; i++ {
				rows[i] = SimpleRow{i}
			}
			model.SetRows(rows)

			got, _ := model.Update(tc.msg)
			wantEq(t, tc.wantCursor, got.Cursor())
		})
	}
}

func wantEq(t *testing.T, want, got interface{}) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %#v, got %#v", want, got)
	}
}
