package main

import "fmt"

type afg struct {
	fg *fg

	srcfn  *fn
	sinkfn *fn

	sfns  []*fn
	rtfns []*fn
	tfns  []*fn

	res *fgr
}

/* MARK: preparation */

func (a *afg) prepFns() error {
	var err error

	a.srcfn = &fn{t: src}
	a.sinkfn = &fn{t: sink}

	a.sfns, err = makeStudentFns()
	if err != nil {
		return err
	}

	a.rtfns, err = makeRTUFns()
	if err != nil {
		return err
	}

	a.tfns, err = makeTutFns()
	if err != nil {
		return err
	}

	return nil
}

func (a *afg) prepFg() error {
	everyFn := append(a.sfns, a.rtfns...)
	everyFn = append(everyFn, a.tfns...)
	everyFn = append(everyFn, a.srcfn, a.sinkfn)

	a.fg = newFg(everyFn)

	for _, sfn := range a.sfns {
		a.fg.addBothCost(a.srcfn, sfn, 1)
	}

	for _, sfn := range a.sfns {
		for _, rtfn := range a.rtfns {
			if sfn.st.open(rtfn.rsec.time) && sfn.st.open(rtfn.tsec.time) {
				a.fg.addBothCost(sfn, rtfn, 1)
			}
		}

	}

	for _, rtfn := range a.rtfns {
		for _, tfn := range a.tfns {
			if rtfn.tsec == tfn.tsec {
				a.fg.addBothCost(rtfn, tfn, int(^uint(0)>>1))
			}
		}
	}

	tutSize := len(a.sfns)/len(a.tfns) + allowedTutOvf
	fmt.Println("ideal tutorial size", tutSize)
	for _, tfn := range a.tfns {
		a.fg.addBothCost(tfn, a.sinkfn, tutSize)
	}

	return nil
}

func (a *afg) prepare(u []*fn) error {
	err := a.prepFns()
	if err != nil {
		return err
	}

	err = a.prepFg()
	if err != nil {
		return err
	}

	return nil
}

func (a *afg) execute() error {
	a.res = a.fg.runMaxFlow(a.srcfn, a.sinkfn)
	if a.res.flow != len(a.sfns) {
		return fmt.Errorf("flow: %d, expected: %d", a.res.flow, len(a.sfns))
	}
	return nil
}

func (a *afg) export() ([]*fn, error) {
	gotrbest := 0
	gotrfos := 0
	gott := 0

	for _, sfn := range a.sfns {
		for _, rtfn := range a.rtfns {
			f := a.res.flowAlong(sfn, rtfn)
			if f > 0 {
				sfn.rsec = rtfn.rsec
				for i, favrsec := range sfn.st.rp {
					if favrsec == rtfn.rsec.time {
						if i == 0 {
							gotrbest++
						}
						gotrfos++
						break
					}
				}
				sfn.tsec = rtfn.tsec
				if sfn.tsec.time == sfn.st.tp {
					gott++
				}
				break
			}
		}
	}

	for _, tfn := range a.tfns {
		f := a.res.flowAlong(tfn, a.sinkfn)
		fmt.Printf("tfn: %s, flow: %d\n", tfn.tsec, f)
	}

	total := len(a.sfns)
	fmt.Println("Ideal recitation fraction: ", float64(gotrbest)/float64(total))
	fmt.Println("Top two recitation fraction: ", float64(gotrfos)/float64(total))
	fmt.Println("Ideal tutorial fraction: ", float64(gott)/float64(total))
	return a.sfns, nil
}
