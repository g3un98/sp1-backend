package main

type account struct {
	Id         string     `json:"id"`
	Pw         string     `json:"pw"`
	Payment    payment    `json:"payment"`
	Membership membership `json:"membership"`
}

func getAccount(ott, id, pw string) (*account, error) {
    switch ott {
    case "Netflix":
        return getNetflixAccount(id, pw)
    case "Wavve":
        return getWavveAccount(id, pw)
    default:
        return &account{
            Id: id,
            Pw: pw,
            Payment: payment{},
            Membership: membership{},
        }, nil
    }
}
