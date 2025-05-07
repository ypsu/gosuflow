# gosuflow: GO StrUct FLOW runner

This package executes linear workflows specially defined via a struct:

    type workflow struct {
      somevar1 type1
      somevar2 type2
      ...
      Step1Section struct{}
      step1var type3
      ...

      Step2Section struct{}
      step2var type4
      ...

      StepNSection struct{}
      stepNvar type5
    }

    func (wf *workflow) Step1(ctx context.Context) error { ... }
    func (wf *workflow) Step2(ctx context.Context) error { ... }
    ...
    func (wf *workflow) StepN(ctx context.Context) error { ... }

    ...
    wf := &workflow{somevar1: ..., somevar2: ...}
    if err := gosuflow.Run(ctx, wf); err != nil {
      ...
    }

gosuflow.Run executes all exported methods the struct has defined.
Each method MethodName must have a corresponding MethodNameSection field in the struct.
gosuflow executes the methods in the order the sections appear in the struct.
By convention each method can only access the fields defined above its section and must initialize the fields belonging to the section.
gosuflow bails out after the first error.

Structuring a long workflow into a gosuflow struct makes the workflow's flow and data much clearer for humans.
gosuflow.Run just fills the boilerplate of calling the functions in the right order and handling the errors.
See the example how a table rendering workflow can be summarized into a gosuflow struct.

See https://pkg.go.dev/github.com/ypsu/gosuflow for the package documentation.
