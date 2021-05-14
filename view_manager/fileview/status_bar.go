package fileview

func (sb *StatusBar) Init(x int) {
        sb.Time.Width  = 8                              //XX:XX:XX
        sb.Time.Origin = (x - 1) - sb.Time.Width        //right aligned

        sb.Status.Origin = 0
        sb.Status.Width  = x - sb.Time.Width

        sb.Elapsed = 21                 //initial value to force clock update
}




