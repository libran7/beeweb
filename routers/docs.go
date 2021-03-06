// Copyright 2013 Beego Web authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package routers

import (
	"net/http"

	"github.com/astaxie/beego/context"
	"github.com/beego/i18n"

	"github.com/beego/beeweb/models"
)

// DocsRouter serves about page.
type DocsRouter struct {
	baseRouter
}

// Get implemented Get method for DocsRouter.
func (this *DocsRouter) Get() {
	this.Data["IsDocs"] = true
	this.TplNames = "docs.html"

	dRoot := models.GetDocByLocale(this.Lang)

	if dRoot == nil {
		this.Abort("404")
		return
	}

	link := this.GetString(":all")

	var doc *models.DocNode
	if len(link) == 0 {
		if dRoot.Doc.HasContent() {
			doc = dRoot.Doc
		} else {
			this.Redirect("/docs/intro/", 302)
			return
		}
	} else {
		doc, _ = dRoot.GetNodeByLink(link)
		if doc == nil {
			doc, _ = dRoot.GetNodeByLink(link + "/")
			if doc != nil {
				this.Redirect("/docs/"+link+"/", 301)
				return
			}
		}
	}

	if doc == nil {
		this.Abort("404")
		return
	}

	this.Data["DocRoot"] = dRoot
	this.Data["Doc"] = doc
	this.Data["Title"] = doc.Name
	this.Data["Data"] = doc.GetContent()
}

func DocsStatic(ctx *context.Context) {
	uri := ctx.Input.Params[":all"]
	if len(uri) > 0 {
		lang := ctx.GetCookie("lang")
		if !i18n.IsExist(lang) {
			lang = "en-US"
		}
		http.ServeFile(ctx.ResponseWriter, ctx.Request, "docs/"+lang+"/"+"images/"+uri)
	}
}
