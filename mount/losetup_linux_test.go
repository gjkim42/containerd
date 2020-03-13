// +build linux

/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package mount

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/containerd/continuity/testutil"
)

func TestSetupLoop(t *testing.T) {
	testutil.RequiresRoot(t)
	const randomdata = "randomdata"

	/* Non-existing loop */
	backingFile := "setup-loop-test-no-such-file"
	_, err := setupLoop(backingFile, LoopParams{})
	if err == nil {
		t.Fatalf("setupLoop with non-existing file should fail")
	}

	f, err := ioutil.TempFile("", "losetup")
	if err != nil {
		t.Fatal(err)
	}
	if err = f.Truncate(512); err != nil {
		t.Fatal(err)
	}
	backingFile = f.Name()
	f.Close()
	defer func() {
		if err := os.Remove(backingFile); err != nil {
			t.Fatal(err)
		}
	}()

	/* RO loop */
	f, err = setupLoop(backingFile, LoopParams{Readonly: true, Autoclear: true})
	if err != nil {
		t.Fatal(err)
	}
	ff, err := os.OpenFile(f.Name(), os.O_RDWR, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = ff.Write([]byte(randomdata)); err == nil {
		t.Fatalf("writing to readonly loop device should fail")
	}
	if err = ff.Close(); err != nil {
		t.Fatal(err)
	}
	if err = f.Close(); err != nil {
		t.Fatal(err)
	}

	/* RW loop */
	f, err = setupLoop(backingFile, LoopParams{Autoclear: true})
	if err != nil {
		t.Fatal(err)
	}
	ff, err = os.OpenFile(f.Name(), os.O_RDWR, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = ff.Write([]byte(randomdata)); err != nil {
		t.Fatal(err)
	}
	if err = ff.Close(); err != nil {
		t.Fatal(err)
	}
	if err = f.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestAttachDetachLoopDevice(t *testing.T) {
	testutil.RequiresRoot(t)
	f, err := ioutil.TempFile("", "losetup")
	if err != nil {
		t.Fatal(err)
	}
	if err = f.Truncate(512); err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer func() {
		if err := os.Remove(f.Name()); err != nil {
			t.Fatal(err)
		}
	}()

	dev, err := AttachLoopDevice(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if err = DetachLoopDevice(dev); err != nil {
		t.Fatal(err)
	}
}
