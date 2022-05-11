package template

import (
	"github.com/Masterminds/sprig/v3"
	"github.com/clevyr/go-yampl/internal/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"strings"
	"text/template"
)

var FuncMap = sprig.TxtFuncMap()

func init() {
	FuncMap["repo"] = DockerRepo
	FuncMap["tag"] = DockerTag
}

func LineComment(conf config.Config, node *yaml.Node) error {
	if node.LineComment != "" && strings.HasPrefix(node.LineComment, conf.Prefix) {
		tmplSrc := strings.TrimSpace(node.LineComment[len(conf.Prefix):])
		tmpl, err := template.New("").
			Funcs(FuncMap).
			Delims(conf.LeftDelim, conf.RightDelim).
			Option("missingkey=error").
			Parse(tmplSrc)
		if err != nil {
			return err
		}

		if conf.Values != nil {
			conf.Values["Value"] = node.Value
		}

		var buf strings.Builder
		if err = tmpl.Execute(&buf, conf.Values); err != nil {
			if !conf.Strict {
				log.WithError(err).Warn("skipping value due to template error")
				return nil
			}
			return err
		}

		if buf.String() != node.Value {
			log.WithFields(log.Fields{
				"tmpl": tmplSrc,
				"from": node.Value,
				"to":   buf.String(),
			}).Debug("updating value")
			node.Value = buf.String()
		}
	}
	return nil
}
