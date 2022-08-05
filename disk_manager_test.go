package bplustree

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestDiskManager_CreateOrOpenFile(t *testing.T) {
	tests := []struct {
		name     string
		arg      string
		want     string
		wantErr  bool
		isDelete bool
	}{
		{
			name:     "Create",
			arg:      "./test/heap_file_create.btr",
			want:     "./test/heap_file_create.btr",
			wantErr:  false,
			isDelete: true,
		},
		{
			name:     "Open",
			arg:      "./test/test_create.btr",
			want:     "./test/test_create.btr",
			wantErr:  false,
			isDelete: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateOrOpenFile(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("failed to CreateOrOpenFile, want: %v, got: %v", tt.wantErr, err)
			}
			if tt.want != got.Name() {
				t.Errorf("%s, want: %s, got: %s", tt.name, tt.want, got.Name())
			}

			// Cleanup
			if tt.isDelete {
				err := os.Remove(tt.arg)
				if err != nil {
					t.Errorf("failed to os.Remove, err: %v", err)
				}
			}
		})
	}
}

func TestDiskManager_GetFileMetadata(t *testing.T) {
	t.Run("Get Metadata", func(t *testing.T) {
		path := "./test/test_create.btr"
		file, err := CreateOrOpenFile(path)
		if err != nil {
			t.Errorf("failed to CreateOrOpenFile, err: %v", err)
		}

		metadata, err := GetFileMetadata(file)
		if err != nil {
			t.Errorf("failed to GetFileMetadata, err: %v", err)
		}

		if metadata == nil {
			t.Errorf("metadata is nil")
		}

		if "test_create.btr" != metadata.Name() {
			t.Errorf("wrong name, %s", metadata.Name())
		}

		if 21 != metadata.Size() {
			t.Errorf("wrong size, %d", metadata.Size())
		}
	})
}

func new4BytetDiskManager(path string) (*DiskManager, error) {
	f, err := CreateOrOpenFile(path)
	if err != nil {
		return nil, err
	}

	metadata, err := GetFileMetadata(f)
	if err != nil {
		return nil, err
	}
	nextPageID := metadata.Size() / 4

	return &DiskManager{
		file:       f,
		nextPageID: nextPageID,
		pageSize:   4,
	}, nil
}

func TestDiskManager_Read(t *testing.T) {
	tests := []struct {
		name    string
		arg     int64
		want    string
		wantErr bool
	}{
		{
			name:    "PageID:0, Data:aaaa",
			arg:     0,
			want:    "aaaa",
			wantErr: false,
		},
		{
			name:    "PageID:2, Data:cccc",
			arg:     2,
			want:    "cccc",
			wantErr: false,
		},
	}

	path := "./test/test_read.btr"
	disk, err := new4BytetDiskManager(path)
	if err != nil {
		t.Errorf("failed to NewDiskManager, err: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := disk.Read(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("failed to Read, want: %v, err: %v", tt.wantErr, err)
			}

			if tt.want != string(got) {
				t.Errorf("%s, want: %s, got: %s", tt.name, tt.want, string(got))
			}
		})
	}
}

func TestDiskManager_Write(t *testing.T) {
	type input struct {
		pageID int64
		data   []byte
	}
	tests := []struct {
		name    string
		input   input
		want    string
		wantErr bool
	}{
		{
			name: "PageID:0, Data:aaaa",
			input: input{
				pageID: 0,
				data:   []byte(`aaaa`),
			},
			want:    "aaaa",
			wantErr: false,
		},
		{
			name: "PageID:1, Data:aaaabbbb",
			input: input{
				pageID: 1,
				data:   []byte(`bbbb`),
			},
			want:    "aaaabbbb",
			wantErr: false,
		},
	}

	// Setup
	path := "./test/test_write.btr"
	disk, err := new4BytetDiskManager(path)
	if err != nil {
		t.Errorf("failed to NewDiskManager, err: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := disk.Write(tt.input.pageID, tt.input.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("failed to Read, want: %v, err: %v", tt.wantErr, err)
			}
			got, err := ioutil.ReadFile(path)
			if err != nil {
				t.Errorf("failed to ReadAl, err: %v", err)
			}
			if tt.want != string(got) {
				t.Errorf("%s, want: %s, got: %s", tt.name, tt.want, string(got))
			}
		})
	}

	// Cleanup
	err = os.Remove(path)
	if err != nil {
		t.Errorf("failed to os.Remove, err: %v", err)
	}
}
