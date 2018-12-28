package stores

import (
	"github.com/d1slike/go-sched/internal"
	"github.com/d1slike/go-sched/triggers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const (
	sName = "test"
)

func TestInMemoryStore_Insert(t *testing.T) {
	store := NewInMemoryStore()
	Convey("Test insertion job and triggers", t, func() {
		Convey("insert job", func() {
			err := store.InsertJob(sName, &internal.Job{Jkey: "job1", JjType: "type1"})
			So(err, ShouldBeNil)
		})

		Convey("insert trigger", func() {
			err := store.InsertTrigger(sName, &internal.Trigger{Tkey: "t1", TjobKey: "job1"})
			So(err, ShouldBeNil)
		})

		Convey("must return err if insert job with same key", func() {
			err := store.InsertJob(sName, &internal.Job{Jkey: "job1", JjType: "type1"})
			So(err, ShouldEqual, ErrJobAlreadyExists)
		})

		Convey("must return err if insert trigger with same key", func() {
			err := store.InsertTrigger(sName, &internal.Trigger{Tkey: "t1", TjobKey: "job1"})
			So(err, ShouldEqual, ErrTriggerAlreadyExists)
		})

		Convey("must return job", func() {
			j, err := store.GetJob(sName, "job1")
			So(err, ShouldBeNil)
			So(j, ShouldNotBeNil)
		})

		Convey("must return trigger", func() {
			t, err := store.GetTrigger(sName, "t1")
			So(err, ShouldBeNil)
			So(t, ShouldNotBeNil)
		})
	})
}

func TestInMemoryStore_Update(t *testing.T) {
	store := NewInMemoryStore()
	Convey("Test updating", t, func() {
		Convey("must return err if no job in store", func() {
			err := store.UpdateJob(sName, &internal.Job{Jkey: "j1"})
			So(err, ShouldEqual, ErrJobNotFound)
		})

		Convey("must return err inf not trigger in store", func() {
			err := store.UpdateTrigger(sName, &internal.Trigger{Tkey: "t1"})
			So(err, ShouldEqual, ErrTriggerNotFound)
		})

		Convey("insert job", func() {
			err := store.InsertJob(sName, &internal.Job{Jkey: "j1", JjType: "type1"})
			So(err, ShouldBeNil)
		})

		Convey("insert trigger", func() {
			err := store.InsertTrigger(sName, &internal.Trigger{Tkey: "t1", Trepeats: 1})
			So(err, ShouldBeNil)
		})

		Convey("update existing job", func() {
			err := store.UpdateJob(sName, &internal.Job{Jkey: "j1", JjType: "type2"})
			So(err, ShouldBeNil)
		})

		Convey("update existing trigger", func() {
			err := store.UpdateTrigger(sName, &internal.Trigger{Tkey: "t1", Trepeats: 2})
			So(err, ShouldBeNil)
		})

		Convey("must return updated job", func() {
			j, err := store.GetJob(sName, "j1")
			So(err, ShouldBeNil)
			So(j.Type(), ShouldEqual, "type2")
		})

		Convey("must return updated trigger", func() {
			j, err := store.GetTrigger(sName, "t1")
			So(err, ShouldBeNil)
			So(j.Repeats(), ShouldEqual, 2)
		})
	})
}

func TestInMemoryStore_Delete(t *testing.T) {
	store := NewInMemoryStore()
	Convey("Test delete", t, func() {
		Convey("must return false if no jobs was deleted", func() {
			ok, err := store.DeleteJob(sName, "j1")
			So(err, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})

		Convey("must return false if no triggers was deleted", func() {
			ok, err := store.DeleteTrigger(sName, "t1")
			So(err, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})

		Convey("insert job", func() {
			err := store.InsertJob(sName, &internal.Job{Jkey: "j1", JjType: "type1"})
			So(err, ShouldBeNil)
		})

		Convey("insert trigger", func() {
			err := store.InsertTrigger(sName, &internal.Trigger{Tkey: "t1", Trepeats: 1})
			So(err, ShouldBeNil)
		})

		Convey("must return true if job was deleted", func() {
			ok, err := store.DeleteJob(sName, "j1")
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)
		})

		Convey("must return true if trigger was deleted", func() {
			ok, err := store.DeleteTrigger(sName, "t1")
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)
		})
	})
}

func TestInMemoryStore_Get(t *testing.T) {
	store := NewInMemoryStore()

	job := internal.Job{Jkey: "j1", JjType: "type1"}
	trigger := internal.Trigger{Tkey: "t1", Trepeats: 1}

	Convey("Test getting", t, func() {
		Convey("must return nil if no job in store", func() {
			j, err := store.GetJob(sName, "j1")
			So(err, ShouldBeNil)
			So(j, ShouldBeNil)
		})

		Convey("must return nil if no triggers in store", func() {
			j, err := store.GetTrigger(sName, "t1")
			So(err, ShouldBeNil)
			So(j, ShouldBeNil)
		})

		Convey("insert job", func() {
			err := store.InsertJob(sName, &job)
			So(err, ShouldBeNil)
		})

		Convey("insert trigger", func() {
			err := store.InsertTrigger(sName, &trigger)
			So(err, ShouldBeNil)
		})

		Convey("jobs are isolated between schedulers", func() {
			j, err := store.GetJob(sName+"1", "j1")
			So(err, ShouldBeNil)
			So(j, ShouldBeNil)
		})

		Convey("triggers are isolated between schedulers", func() {
			t, err := store.GetTrigger(sName+"1", "t1")
			So(err, ShouldBeNil)
			So(t, ShouldBeNil)
		})

		Convey("must return exactly equal job", func() {
			j, err := store.GetJob(sName, "j1")
			So(err, ShouldBeNil)
			So(j, ShouldNotBeNil)
			So(j, ShouldResemble, &job)
		})

		Convey("must return exactly equal trigger", func() {
			t, err := store.GetTrigger(sName, "t1")
			So(err, ShouldBeNil)
			So(t, ShouldResemble, &trigger)
		})

		Convey("insert one more job", func() {
			err := store.InsertJob(sName, &internal.Job{Jkey: "j2"})
			So(err, ShouldBeNil)
			err = store.InsertJob(sName+"1", &internal.Job{Jkey: "j2"})
			So(err, ShouldBeNil)
		})

		Convey("insert one more trigger", func() {
			err := store.InsertTrigger(sName, &internal.Trigger{Tkey: "t2"})
			So(err, ShouldBeNil)
			err = store.InsertTrigger(sName+"1", &internal.Trigger{Tkey: "t2"})
			So(err, ShouldBeNil)
		})

		Convey("total jobs count must equal 2", func() {
			jobs, err := store.GetJobs(sName)
			So(err, ShouldBeNil)
			So(jobs, ShouldHaveLength, 2)
		})

		Convey("total triggers count must equal 2", func() {
			triggers, err := store.GetTriggers(sName)
			So(err, ShouldBeNil)
			So(triggers, ShouldHaveLength, 2)
		})
	})
}

func TestInMemoryStore_DeleteExhaustedTriggers(t *testing.T) {
	store := NewInMemoryStore()

	Convey("Test deleting exhausted triggers", t, func() {
		Convey("insert triggers", func() {
			err := store.InsertTrigger(sName, &internal.Trigger{Tkey: "t1", Tstate: triggers.StateAcquired})
			So(err, ShouldBeNil)
			err = store.InsertTrigger(sName, &internal.Trigger{Tkey: "t2", Tstate: triggers.StateExhausted})
			So(err, ShouldBeNil)
			err = store.InsertTrigger(sName, &internal.Trigger{Tkey: "t3", Tstate: triggers.StateScheduled})
			So(err, ShouldBeNil)
			err = store.InsertTrigger(sName+"1", &internal.Trigger{Tkey: "t4", Tstate: triggers.StateExhausted})
			So(err, ShouldBeNil)
		})

		Convey("must delete 1 trigger", func() {
			count, err := store.DeleteExhaustedTriggers(sName)
			So(err, ShouldBeNil)
			So(count, ShouldEqual, 1)

			triggers, err := store.GetTriggers(sName)
			So(err, ShouldBeNil)
			So(triggers, ShouldHaveLength, 2)

			triggers, err = store.GetTriggers(sName + "1")
			So(err, ShouldBeNil)
			So(triggers, ShouldHaveLength, 1)

			t, err := store.GetTrigger(sName, "t2")
			So(err, ShouldBeNil)
			So(t, ShouldBeNil)
		})
	})
}

func TestInMemoryStore_AcquireTriggers(t *testing.T) {
	store := NewInMemoryStore()

	Convey("Test trigger acquiring", t, func() {
		Convey("insert triggers", func() {
			err := store.InsertTrigger(sName, &internal.Trigger{Tkey: "t1", Tstate: triggers.StateAcquired})
			So(err, ShouldBeNil)
			err = store.InsertTrigger(sName, &internal.Trigger{Tkey: "t2", Tstate: triggers.StateExhausted})
			So(err, ShouldBeNil)
			err = store.InsertTrigger(sName, &internal.Trigger{Tkey: "t3", Tstate: triggers.StateScheduled})
			So(err, ShouldBeNil)
			err = store.InsertTrigger(sName, &internal.Trigger{Tkey: "t4", Tstate: triggers.StateScheduled})
			So(err, ShouldBeNil)
			err = store.InsertTrigger(sName+"1", &internal.Trigger{Tkey: "t5", Tstate: triggers.StateScheduled})
			So(err, ShouldBeNil)
		})

		Convey("must return 2 triggers", func() {
			arr, err := store.AcquireTriggers(sName)
			So(err, ShouldBeNil)
			So(arr, ShouldHaveLength, 2)
			checkStatus := true
			for _, v := range arr {
				if v.State() != triggers.StateAcquired {
					checkStatus = false
					break
				}
			}
			So(checkStatus, ShouldBeTrue)

			t, err := store.GetTrigger(sName, "t3")
			So(err, ShouldBeNil)
			So(t, ShouldNotBeNil)
			So(t.State(), ShouldEqual, triggers.StateAcquired)
		})

		Convey("must return empty array", func() {
			arr, err := store.AcquireTriggers(sName)
			So(err, ShouldBeNil)
			So(arr, ShouldBeEmpty)
		})
	})
}

func TestInMemoryStore_DeleteTriggersByJobKey(t *testing.T) {
	store := NewInMemoryStore()

	Convey("Test delete related triggers", t, func() {
		Convey("insert entities", func() {
			err := store.InsertJob(sName, &internal.Job{Jkey: "job1"})
			So(err, ShouldBeNil)
			err = store.InsertTrigger(sName, &internal.Trigger{Tkey: "t1", TjobKey: "job1"})
			So(err, ShouldBeNil)
			err = store.InsertTrigger(sName, &internal.Trigger{Tkey: "t2", TjobKey: "job1"})
			So(err, ShouldBeNil)
			err = store.InsertTrigger(sName, &internal.Trigger{Tkey: "t3", TjobKey: "job3"})
			So(err, ShouldBeNil)
			err = store.InsertTrigger(sName+"1", &internal.Trigger{Tkey: "t1", TjobKey: "job1"})
			So(err, ShouldBeNil)
		})

		Convey("must delete 2 triggers", func() {
			arr, err := store.DeleteTriggersByJobKey(sName, "job1")
			So(err, ShouldBeNil)
			So(arr, ShouldResemble, []string{"t1", "t2"})

			t, err := store.GetTrigger(sName, "t1")
			So(err, ShouldBeNil)
			So(t, ShouldBeNil)
		})
	})
}
