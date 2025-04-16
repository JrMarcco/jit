package copier

import (
	"reflect"
	"slices"
	"time"

	"maps"

	"github.com/JrMarcco/easy-kit/bean/opt"
	"github.com/JrMarcco/easy-kit/set"
)

type RefCopier[S any, D any] struct {
	root        fieldNode
	atomicTypes []reflect.Type

	defaultConf copyConf
}

func NewRefCopier[S any, D any](opts ...opt.Opt[copyConf]) (*RefCopier[S, D], error) {
	srcTyp := reflect.TypeOf(new(S)).Elem()
	dstTyp := reflect.TypeOf(new(D)).Elem()

	if srcTyp.Kind() != reflect.Struct {
		return nil, errInvalidType("struct", srcTyp)
	}

	if dstTyp.Kind() != reflect.Struct {
		return nil, errInvalidType("struct", dstTyp)
	}

	root := fieldNode{
		fields: []fieldNode{},
	}

	copier := &RefCopier[S, D]{
		root: root,
		atomicTypes: []reflect.Type{
			reflect.TypeOf(time.Time{}),
		},
	}

	if err := copier.createFieldNode(srcTyp, dstTyp, &root); err != nil {
		return nil, err
	}

	copier.root = root

	cc := newCopyConf()
	opt.Apply(&cc, opts...)

	copier.defaultConf = cc
	return copier, nil
}

func (rc *RefCopier[S, D]) createFieldNode(srcTyp, dstTyp reflect.Type, root *fieldNode) error {
	srcMap := map[string]int{}
	for i := range srcTyp.NumField() {
		fd := srcTyp.Field(i)
		if fd.IsExported() {
			srcMap[fd.Name] = i
		}
	}

	for dstIdx := range dstTyp.NumField() {
		dstFd := dstTyp.Field(dstIdx)

		if !dstFd.IsExported() {
			continue
		}

		if srcIdx, ok := srcMap[dstFd.Name]; ok {
			srcFd := srcTyp.Field(srcIdx)

			if srcFd.Type.Kind() == reflect.Pointer && srcFd.Type.Elem().Kind() == reflect.Pointer {
				// pointer to pointer
				return errPtrToPtr(srcFd.Name)
			}

			if dstFd.Type.Kind() == reflect.Pointer && dstFd.Type.Elem().Kind() == reflect.Pointer {
				// pointer to pointer
				return errPtrToPtr(dstFd.Name)
			}

			node := fieldNode{
				name:   dstFd.Name,
				sIndex: srcIdx,
				dIndex: dstIdx,
				fields: []fieldNode{},
			}

			srcFdTyp := srcFd.Type
			dstFdTyp := dstFd.Type

			if srcFdTyp.Kind() == reflect.Pointer {
				srcFdTyp = srcFdTyp.Elem()
			}

			if dstFdTyp.Kind() == reflect.Pointer {
				dstFdTyp = dstFdTyp.Elem()
			}

			if isBuiltinType(srcFdTyp.Kind()) {
				// builtin type, node is leaf node
			} else if rc.isAtomicType(srcFdTyp) {
				// is atomic type, node is leaf node
			} else if srcFdTyp.Kind() == reflect.Struct {
				// is struct type
				if err := rc.createFieldNode(srcFdTyp, dstFdTyp, &node); err != nil {
					return err
				}
			} else {
				// is not builtin type, not struct type, not atomic type
				// can not copy it, skip
				continue
			}

			// add node to root.fields
			root.fields = append(root.fields, node)
		}
	}
	return nil
}

func (rc *RefCopier[S, D]) isAtomicType(typ reflect.Type) bool {
	return slices.Contains(rc.atomicTypes, typ)
}

func (rc *RefCopier[S, D]) Copy(src *S, opts ...opt.Opt[copyConf]) (*D, error) {
	dst := new(D)
	err := rc.CopyTo(src, dst, opts...)
	return dst, err
}

func (rc *RefCopier[S, D]) CopyTo(src *S, dst *D, opts ...opt.Opt[copyConf]) error {
	if len(rc.root.fields) == 0 {
		return nil
	}

	cc := rc.defaultCopyConf()
	opt.Apply(&cc, opts...)
	return rc.copyTree(src, dst, cc)
}

func (rc *RefCopier[S, D]) defaultCopyConf() copyConf {
	cc := newCopyConf()

	if rc.defaultConf.ignoreFds != nil {
		ignoreFds := set.NewMapSet[string](rc.defaultConf.ignoreFds.Size())

		for _, fd := range rc.defaultConf.ignoreFds.Elems() {
			ignoreFds.Add(fd)
		}

		cc.ignoreFds = ignoreFds
	}

	if cc.covertFds == nil {
		cc.covertFds = make(map[string]convertFunc, len(rc.defaultConf.covertFds))
	}

	maps.Copy(cc.covertFds, rc.defaultConf.covertFds)

	return cc
}

func (rc *RefCopier[S, D]) copyTree(src *S, dst *D, cc copyConf) error {
	srcTyp := reflect.TypeOf(src)
	srcVal := reflect.ValueOf(src)

	dstTyp := reflect.TypeOf(dst)
	dstVal := reflect.ValueOf(dst)

	return rc.copyNode(srcTyp, srcVal, dstTyp, dstVal, &rc.root, cc)
}

func (rc *RefCopier[S, D]) copyNode(srcTyp reflect.Type, srcVal reflect.Value, dstTyp reflect.Type, dstVal reflect.Value, root *fieldNode, cc copyConf) error {
	oriSrcVal := srcVal
	oriDstVal := dstVal

	if srcVal.Kind() == reflect.Pointer {
		if srcVal.IsNil() {
			return nil
		}
		srcVal = srcVal.Elem()
		srcTyp = srcTyp.Elem()
	}

	if dstVal.Kind() == reflect.Pointer {
		if dstVal.IsNil() {
			dstVal.Set(reflect.New(dstTyp.Elem()))
		}

		dstVal = dstVal.Elem()
		dstTyp = dstTyp.Elem()
	}

	if len(root.fields) == 0 {
		fdName := root.name
		if !dstVal.CanSet() {
			return nil
		}

		convertFunc, ok := cc.covertFds[fdName]
		if !ok {
			if srcTyp != dstTyp {
				return errFieldTypeMismatch(fdName, srcTyp, dstTyp)
			}

			if srcVal.IsZero() {
				return nil
			}

			dstVal.Set(srcVal)
			return nil
		}

		if !oriDstVal.CanSet() {
			return nil
		}

		srcConverted, err := convertFunc(oriSrcVal.Interface())
		if err != nil {
			return err
		}

		srcConvTyp := reflect.TypeOf(srcConverted)
		srcConvVal := reflect.ValueOf(srcConverted)

		if srcConvTyp != oriDstVal.Type() {
			return errFieldTypeMismatch(fdName, srcConvTyp, oriDstVal.Type())
		}

		oriDstVal.Set(srcConvVal)

		return nil
	}

	for _, field := range root.fields {
		if cc.InIngore(field.name) {
			continue
		}

		srcFdTyp := srcTyp.Field(field.sIndex)
		srcFdVal := srcVal.Field(field.sIndex)

		dstFdTyp := dstTyp.Field(field.dIndex)
		dstFdVal := dstVal.Field(field.dIndex)

		if err := rc.copyNode(srcFdTyp.Type, srcFdVal, dstFdTyp.Type, dstFdVal, &field, cc); err != nil {
			return err
		}
	}

	return nil
}

type fieldNode struct {
	name   string
	fields []fieldNode
	sIndex int // source index
	dIndex int // destination index
}

func isBuiltinType(kind reflect.Kind) bool {
	switch kind {
	case
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
		reflect.String,
		reflect.Slice,
		reflect.Map,
		reflect.Chan,
		reflect.Array:
		return true
	default:
		return false
	}
}
