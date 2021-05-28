module github.com/rsmithsa/digital-ocean-ddns

go 1.16

require golang.org/x/net v0.0.0-20210525063256-abc453219eb5
require internal/doapiv2 v1.0.0
replace internal/doapiv2 => ./internal/doapiv2