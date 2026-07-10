package main

import (
	"fmt"
	"strings"
)

// General use function to format strings for PlantUML
func getPumlFormat(s string) string {
	r := strings.NewReplacer("-", "_", ".", "_", " ", "_", "/", "_", "(", "_", ")", "_")
	return strings.ToLower(r.Replace(s))
}

// VLAN represents a network VLAN with a unique identifier.
type vlan struct {
	vid int
}

// Termination represents one end of a network cable.
type termination struct {
	device_name string
	name        string
}

func (t termination) getPuml() string {
	return fmt.Sprintf("%s::%s", getPumlFormat(t.device_name), getPumlFormat(t.name))
}

// Cable represents a network cable connecting two terminations.
type cable struct {
	status         string
	a_terminations []termination
	b_terminations []termination
}

func (c cable) isConnected() bool {
	return c.status == "connected" || c.status == "decommissioning"
}

func (c cable) getConnPumls() []string {
	if c.isConnected() {
		var pumls []string
		for _, t := range c.a_terminations {
			for _, t2 := range c.b_terminations {
				var s = fmt.Sprintf("%s <--> %s", t.getPuml(), t2.getPuml())
				pumls = append(pumls, s)
			}
		}
		return pumls
	}
	return nil
}

// Connector represents a network connector that can be connected and has a PlantUML representation.
type connector interface {
	isConnected() bool
	getPuml() (string, []string)
}

// Port represents a network port that can be connected to a cable.
type port struct {
	name  string
	cable *cable
}

func (p *port) isConnected() bool {
	if p.cable != nil {
		return p.cable.isConnected()
	}
	return false
}
func (p *port) getPuml() (string, []string) {
	if p.isConnected() {
		return getPumlFormat(p.name), p.cable.getConnPumls()
	}
	return "", nil
}

// Interface represents a network interface that can be connected to a cable and may have VLAN configurations.
type interf struct {
	name          string
	cable         *cable
	enabled       bool
	untagged_vlan *vlan
	tagged_vlans  []*vlan
}

func (i *interf) isConnected() bool {
	return i.enabled && i.cable != nil && i.cable.isConnected()
}

func (i *interf) getPuml() (string, []string) {
	if i.isConnected() {
		puml := getPumlFormat(i.name)

		puml += " => "
		if i.untagged_vlan != nil {
			puml += fmt.Sprintf("(%d)", i.untagged_vlan.vid)
		}
		if len(i.tagged_vlans) > 0 {
			if i.untagged_vlan != nil {
				puml += ","
			}
			var tagged []string
			for _, v := range i.tagged_vlans {
				tagged = append(tagged, fmt.Sprintf("%d", v.vid))
			}
			puml += strings.Join(tagged, ",")
		}

		return puml, i.cable.getConnPumls()
	}
	return "", nil
}

// Device represents a network device with front ports, rear ports, and interfaces.
type device struct {
	name       string
	frontports []port
	rearports  []port
	interfaces []interf
}

func (d device) getPuml() (string, []string) {
	var pumls []string
	var conn_pumls []string

	if d.interfaces != nil {
		pumls = append(pumls, fmt.Sprintf("map \"%s\" as %s {", d.name, getPumlFormat(d.name)))
		for _, i := range d.interfaces {
			if i.isConnected() {
				p, c := i.getPuml()
				pumls = append(pumls, p)
				conn_pumls = append(conn_pumls, c...)
			}
		}
	} else {
		pumls = append(pumls, fmt.Sprintf("object \"%s\" as %s {", d.name, getPumlFormat(d.name)))
		ports := append(d.frontports, d.rearports...)
		for _, port := range ports {
			if port.isConnected() {
				p, c := port.getPuml()
				pumls = append(pumls, p)
				conn_pumls = append(conn_pumls, c...)
			}
		}
	}
	pumls = append(pumls, "}")
	pumls = append(pumls, "\n")

	return strings.Join(pumls, "\n"), conn_pumls
}

func (ds []device) getPuml() string {

}

func main() {
	fmt.Println("Hello, World!")
}
