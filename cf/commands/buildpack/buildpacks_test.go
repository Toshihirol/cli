package buildpack_test

import (
	testapi "github.com/cloudfoundry/cli/cf/api/fakes"
	"github.com/cloudfoundry/cli/cf/commands/buildpack"
	"github.com/cloudfoundry/cli/cf/models"
	testcmd "github.com/cloudfoundry/cli/testhelpers/commands"
	testreq "github.com/cloudfoundry/cli/testhelpers/requirements"
	testterm "github.com/cloudfoundry/cli/testhelpers/terminal"

	. "github.com/cloudfoundry/cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListBuildpacks", func() {
	var (
		ui                  *testterm.FakeUI
		buildpackRepo       *testapi.FakeBuildpackRepository
		requirementsFactory *testreq.FakeReqFactory
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		buildpackRepo = &testapi.FakeBuildpackRepository{}
		requirementsFactory = &testreq.FakeReqFactory{}
	})

	runCommand := func(args ...string) bool {
		cmd := buildpack.NewListBuildpacks(ui, buildpackRepo)
		return testcmd.RunCommand(cmd, args, requirementsFactory)
	}

	It("should fail with usage when provided any arguments", func() {
		requirementsFactory.LoginSuccess = true
		Expect(runCommand("blahblah")).To(BeFalse())
		Expect(ui.FailedWithUsage).To(BeTrue())
	})

	It("fails requirements when login fails", func() {
		Expect(runCommand()).To(BeFalse())
	})

	Context("when logged in", func() {
		BeforeEach(func() {
			requirementsFactory.LoginSuccess = true
		})

		It("lists buildpacks", func() {
			p1 := 5
			p2 := 10
			p3 := 15
			t := true
			f := false

			buildpackRepo.Buildpacks = []models.Buildpack{
				models.Buildpack{Name: "Buildpack-1", Position: &p1, Enabled: &t, Locked: &f},
				models.Buildpack{Name: "Buildpack-2", Position: &p2, Enabled: &f, Locked: &t},
				models.Buildpack{Name: "Buildpack-3", Position: &p3, Enabled: &t, Locked: &f},
			}

			runCommand()

			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Getting buildpacks"},
				[]string{"buildpack", "position", "enabled"},
				[]string{"Buildpack-1", "5", "true", "false"},
				[]string{"Buildpack-2", "10", "false", "true"},
				[]string{"Buildpack-3", "15", "true", "false"},
			))
		})

		It("tells the user if no build packs exist", func() {
			runCommand()
			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Getting buildpacks"},
				[]string{"No buildpacks found"},
			))
		})
	})

})
