package main

type StringSet struct {

  values []string
}

func (this *StringSet) Add(item string) {

  if(!this.Contains(item)) {
    this.values = append(this.values, item)
  }
}

func (this StringSet) Contains(item string) bool {

  for _, extant := range this.values {

    if(extant == item) {
      return true
    }
  }

  return false
}

func (this StringSet) GetSlice() []string {
  return this.values
}
