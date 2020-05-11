// Marshalling (and conversely unmarshalling) XML in Go is weird and
// dumb and it confuses me. Maybe part of this confusion stems from
// this comment at the end of the XML package:
// https://golang.org/pkg/encoding/xml/
//
// Mapping between XML elements and data structures is inherently
// flawed: an XML element is an order-dependent collection of
// anonymous values, while a data structure is an order-independent
// collection of named values. See package json for a textual
// representation more suitable to data structures.
//
// I think one of the things that I have trouble keeping in mind is
// that EVERY field on a struct will produce an XML node containing
// the value of the field.
package main

import (
	"encoding/xml"
	"fmt"
)

// If I marshal this struct on it's own then it will look like:
// <Name><FirstName>...</FirstName><LastName>...</LastName></Name>
type Name struct {
	FirstName string
	LastName  string
}

// If I marshal this struct on it's own then it will look like:
// <NodeNameForBodyStats xmlns="xml-namespace-thingy"><HeightInFeet
// xmlns="another-ns-thingy-cause-why-not">6</HeightInFeet><weight>160</weight></NodeNameForBodyStats>
type BodyStats struct {
	XMLName xml.Name `xml:"xml-namespace-thingy NodeNameForBodyStats"`
	Height  int      `xml:"another-ns-thingy-cause-why-not HeightInFeet"`
	Weight  int      `xml:"weight"`
}

// If I marshal this struct then it will marshal like:
// <PersonName><person-name><FirstName>Lucas</FirstName><LastName>Groenendaal</LastName></person-name></PersonName>
// Note that the "Name" struct type marshalled DIFFERENTLY than if you
// marshalled it by itself. The xml node name specified at this level
// overrode the one that would get generated by default.
type PersonName struct {
	Name Name `xml:"person-name"`
}

// This will fail to marshal because there is a XML node name
// specified here that is DIFFERENT than the one specified on the
// BodyStats struct.
type PersonBodyStats struct {
	BodyStats BodyStats `xml:"different-name-than-specified-in-BodyStats"`
}

// If I pass in the Name and BodyStats structs above into this one
// then it will marshal as:
//
// <PersonInterface><will-dictate-node-value><FirstName></FirstName><LastName></LastName></will-dictate-node-value><NodeNameForBodyStats
// xmlns="xml-namespace-thingy"><HeightInFeet
// xmlns="another-ns-thingy-cause-why-not">0</HeightInFeet><weight>0</weight></NodeNameForBodyStats></PersonInterface>
//
// Note that the xml node name specified on this level overwrote the
// Name struct because that one did not explicitly set a XML name but
// it did NOT override the BodyStats one because that type did specify
// a XML name.
type PersonInterface struct {
	Name      interface{} `xml:"will-dictate-node-value"`
	BodyStats interface{} `xml:"will-not-dictate-node-value"`
}

// When marshalling embedded information it is truly like you've
// pulled out the fields within each embedded struct and threw them
// into this one. In this particular case, the node encompassing all
// this data will be NodeNameForBodyStats because that is the only
// xml.Name value between the structs PersonEmbed, Name, and
// BodyStats: <NodeNameForBodyStats
// xmlns="xml-namespace-thingy"><FirstName></FirstName><LastName></LastName><HeightInFeet
// xmlns="another-ns-thingy-cause-why-not">0</HeightInFeet><weight>0</weight></NodeNameForBodyStats>
type PersonEmbed struct {
	Name
	BodyStats `xml:"this-will-do-nothing"`
}

// When there are xml.Name values on some of the embedded structs, the
// xml.Name on the struct *doing* the embedding takes prescedence:
// <person-embed-xml-name><FirstName></FirstName><LastName></LastName><HeightInFeet
// xmlns="another-ns-thingy-cause-why-not">0</HeightInFeet><weight>0</weight><field>hey</field></person-embed-xml-name>
type PersonEmbedXMLName struct {
	XMLName xml.Name `xml:"person-embed-xml-name"`
	Name
	BodyStats
	OneMoreField string `xml:"field"`
}

type BloodType struct {
	XMLName  xml.Name `xml:"blood-type"`
	Letter   string
	Positive bool
}

// When embedding multiple structs each with a xml.Name value on it
// the one defined first determines the XML node name which will hold
// this value:
// <blood-type><Letter></Letter><Positive>false</Positive><HeightInFeet
// xmlns="another-ns-thingy-cause-why-not">0</HeightInFeet><weight>0</weight></blood-type>
type PersonEmbedsStructsWithXMLNames struct {
	BloodType
	BodyStats
}

// This just goes to show yet again that an XML tag applied to an
// embedded struct does NOTHING.
// <PersonEmbedNoXMLName><FirstName></FirstName><LastName></LastName></PersonEmbedNoXMLName>
type PersonEmbedNoXMLName struct {
	Name `xml:"even-though-Name-struct-has-no-xml-node-this-will-not-work"`
}

// This is one way to achieve some deep nesting and not need to write
// as many structs: <layers-ma-boy-layers><l1str>first
// layer</l1str><layer2><l2str>second
// layer</l2str><layer3><l3str>third
// layer</l3str></layer3><another-l2-thing></another-l2-thing></layer2></layers-ma-boy-layers>
type Layers1 struct {
	XMLName        xml.Name `xml:"layers-ma-boy-layers"`
	L1Str          string   `xml:"l1str"`
	L2Str          string   `xml:"layer2>l2str"`
	L3Str          string   `xml:"layer2>layer3>l3str"`
	AnotherL2Thing string   `xml:"layer2>another-l2-thing"`
}

type Layer23 struct {
	L3Str string `xml:"l3str"`
}

type Layer22 struct {
	L2Str          string  `xml:"l2str"`
	Layer23        Layer23 `xml:"layer3"`
	AnotherL2Thing string  `xml:"another-l2-thing"`
}

// This will marshal the same as the Layers1 struct but it involves
// creating more structs to do the nesting.
type Layers2 struct {
	XMLName xml.Name `xml:"layers-ma-boy-layers"`
	L1Str   string   `xml:"l1str"`
	Layer22 Layer22  `xml:"layer2"`
}

type Layer33 struct {
	XMLName xml.Name `xml:"layer3"`
	L3Str   string   `xml:"l3str"`
}

type Layer32 struct {
	XMLName        xml.Name `xml:"layer2"`
	L2Str          string   `xml:"l2str"`
	Layer33        Layer33
	AnotherL2Thing string `xml:"another-l2-thing"`
}

// And this is yet another way to marshal the same XML as Layers1 and
// Layers2. This time we specified the XML node names for structs
// strictly on xml.Name values in the struct where in Layers2 we had
// the struct which embed another struct decide what the node name for
// that struct would be. I think the Layer2 would be my preferred
// approach (where the only time you use xml.Name is in the struct
// defining the root of the document) because it makes all other
// structs more flexible in where they can be unmarshalled.
type Layers3 struct {
	XMLName xml.Name `xml:"layers-ma-boy-layers"`
	L1Str   string   `xml:"l1str"`
	Layer32 Layer32
}

func xmlMarshalAndPrint(v interface{}) {
	b, err := xml.Marshal(v)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("%s\n", b)
}

func main() {
	xmlMarshalAndPrint(Name{FirstName: "Lucas", LastName: "Groenendaal"})
	xmlMarshalAndPrint(BodyStats{Height: 6, Weight: 160})
	xmlMarshalAndPrint(PersonName{Name: Name{FirstName: "Lucas", LastName: "Groenendaal"}})
	xmlMarshalAndPrint(PersonBodyStats{BodyStats: BodyStats{}})
	xmlMarshalAndPrint(PersonInterface{Name: Name{}, BodyStats: BodyStats{}})
	xmlMarshalAndPrint(PersonEmbed{Name: Name{}, BodyStats: BodyStats{}})
	xmlMarshalAndPrint(PersonEmbedXMLName{Name: Name{}, BodyStats: BodyStats{}, OneMoreField: "hey"})
	xmlMarshalAndPrint(PersonEmbedsStructsWithXMLNames{BloodType: BloodType{}, BodyStats: BodyStats{}})
	xmlMarshalAndPrint(PersonEmbedNoXMLName{})
	xmlMarshalAndPrint(Layers1{L1Str: "first layer", L2Str: "second layer", L3Str: "third layer"})
	xmlMarshalAndPrint(Layers2{L1Str: "first layer", Layer22: Layer22{L2Str: "second layer", Layer23: Layer23{L3Str: "third layer"}}})
	xmlMarshalAndPrint(Layers3{L1Str: "first layer", Layer32: Layer32{L2Str: "second layer", Layer33: Layer33{L3Str: "third layer"}}})
}