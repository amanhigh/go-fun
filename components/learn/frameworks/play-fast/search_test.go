package play_fast

import (
	"sort"
	"time"

	"github.com/lithammer/fuzzysearch/fuzzy"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gmeasure"
)

var _ = Describe("Search", func() {

	Context("FuzzySearch", func() {

		Context("Basic", func() {
			It("1.1 should match strings with partial input", func() {
				Expect(fuzzy.Match("twl", "cartwheel")).To(BeTrue())
				Expect(fuzzy.Match("cart", "cartwheel")).To(BeTrue())
				Expect(fuzzy.Match("cw", "cartwheel")).To(BeTrue())
			})

			It("1.2 should match with character transpositions and missing characters", func() {
				Expect(fuzzy.Match("ee", "cartwheel")).To(BeTrue())
				Expect(fuzzy.Match("whl", "cartwheel")).To(BeTrue())
			})

			It("1.3 should return false for non-matching cases", func() {
				Expect(fuzzy.Match("cwm", "cartwheel")).To(BeFalse())
				Expect(fuzzy.Match("kitten", "sitting")).To(BeFalse())
			})

			It("1.4 should get Levenshtein distance scores via RankMatch", func() {
				rank := fuzzy.RankMatch("cart", "cartwheel")
				Expect(rank).To(Equal(5))
			})

			It("1.4 should return higher scores for better matches", func() {
				exactRank := fuzzy.RankMatch("cart", "cart")
				partialRank := fuzzy.RankMatch("cart", "cartwheel")
				Expect(exactRank).To(Equal(0))
				Expect(partialRank).To(Equal(5))
			})

			It("1.5 should return -1 for non-matches", func() {
				Expect(fuzzy.RankMatch("kitten", "sitting")).To(Equal(-1))
			})
		})

		Context("Medium", func() {
			var words []string

			BeforeEach(func() {
				words = []string{"cartwheel", "foobar", "wheel", "baz", "cart"}
			})

			It("2.1 should find matches in string slice using Find", func() {
				results := fuzzy.Find("whl", words)
				Expect(results).To(ContainElement("cartwheel"))
				Expect(results).To(ContainElement("wheel"))
				Expect(results).ToNot(ContainElement("foobar"))
			})

			It("2.2 should return ranked results with RankFind sorted by score", func() {
				results := fuzzy.RankFind("whl", words)
				Expect(results).ToNot(BeEmpty())
				sort.Sort(results)
				Expect(results[0].Target).To(Equal("wheel"))
			})

			It("2.3 should handle case-insensitive matching with MatchFold", func() {
				Expect(fuzzy.MatchFold("CART", "cartwheel")).To(BeTrue())
				Expect(fuzzy.MatchFold("Cart", "cartwheel")).To(BeTrue())
			})

			It("2.4 should handle Unicode normalized matching with MatchNormalized", func() {
				Expect(fuzzy.MatchNormalized("café", "cafe")).To(BeTrue())
			})
		})

		Context("Advanced", func() {
			It("3.1 should filter autocomplete suggestions", func() {
				commands := []string{
					"git commit", "git push", "git pull", "git status",
					"git branch", "git checkout", "git merge", "git rebase",
					"docker build", "docker run", "docker ps",
				}
				results := fuzzy.Find("git co", commands)
				Expect(results).To(ContainElement("git commit"))
				Expect(results).To(ContainElement("git checkout"))
			})

			It("3.2 should suggest typo corrections", func() {
				cities := []string{
					"New York", "Los Angeles", "Chicago", "Houston",
					"Phoenix", "Philadelphia", "San Antonio", "San Diego",
				}
				results := fuzzy.Find("Chcago", cities)
				Expect(results).To(ContainElement("Chicago"))
			})

			It("3.3 should support command palette / quick open functionality", func() {
				files := []string{
					"main.go", "handler.go", "manager.go", "repository.go",
					"config.yaml", "README.md", "Makefile", "docker-compose.yml",
				}
				results := fuzzy.Find("mng", files)
				Expect(results).To(ContainElement("manager.go"))
			})

			It("3.4 should handle empty pattern", func() {
				Expect(fuzzy.Match("", "anything")).To(BeTrue())
			})

			It("3.5 should handle empty dataset", func() {
				results := fuzzy.Find("query", []string{})
				Expect(results).To(BeEmpty())
			})

			It("3.6 should handle Unicode and special characters", func() {
				Expect(fuzzy.Match("日本", "日本語")).To(BeTrue())
				Expect(fuzzy.Match("über", "übermensch")).To(BeTrue())
			})

			It("3.7 should demonstrate case sensitivity behavior", func() {
				Expect(fuzzy.Match("CART", "cartwheel")).To(BeFalse())
				Expect(fuzzy.MatchFold("CART", "cartwheel")).To(BeTrue())
			})

			Context("Performance Benchmarks", FlakeAttempts(3), func() {
				// BUG: Convert to Benchmark Test.
				var dataset []string

				BeforeEach(func() {
					dataset = make([]string, 10000)
					for i := range dataset {
						dataset[i] = "item-" + string(rune('a'+i%26)) + string(rune('a'+i%17)) + string(rune('a'+i%13))
					}
				})

				It("3.8 should perform Match operations efficiently", func() {
					experiment := gmeasure.NewExperiment("FuzzySearch Match")
					AddReportEntry(experiment.Name, experiment)
					experiment.SampleDuration("match", func(_ int) {
						fuzzy.Match("abc", "abcdefghij")
					}, gmeasure.SamplingConfig{N: 10000})
					Expect(experiment.GetStats("match").DurationFor(gmeasure.StatMedian)).To(
						BeNumerically("<", 1*time.Microsecond))
				})

				It("3.8 should perform Find on large datasets within acceptable latency", func() {
					experiment := gmeasure.NewExperiment("FuzzySearch Find")
					AddReportEntry(experiment.Name, experiment)
					experiment.SampleDuration("find", func(_ int) {
						fuzzy.Find("abc", dataset)
					}, gmeasure.SamplingConfig{N: 100})
					AddReportEntry("Find Stats", experiment.GetStats("find"))
					Expect(experiment.GetStats("find").DurationFor(gmeasure.StatMedian)).To(
						BeNumerically("<", 10*time.Millisecond))
				})
			})
		})
	})
})
