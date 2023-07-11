package valueobject

import "github.com/jacklv111/aifs/pkg/annotation-template/do"

type AnnotationTemplateDetails struct {
	AnnotationTemplateDo    do.AnnotationTemplateDo
	AnnotationTemplateExtDo do.AnnotationTemplateExtDo
	LabelDoList             []do.LabelDo
}
