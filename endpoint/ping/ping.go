/*
This service provides authentication and authorization via OAuth2.

Copyright (C) 2020 Lars Gaubisch

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package ping

type ping struct {
    svc *smis.Service
}

// Init initialises the ping endpoints
func Init(svc *smis.Service) error {
    endpoint := &ping{svc: svc}
    _, err := svc.RegisterEndpoint("/ping", http.MethodGet, endpoint.pingHandler)

    return err
}

func (p *ping) pingHandler(writer http.ResponseWriter, request *http.Request) {
    log := p.svc.NewLogForRequestID(request.Context())

    _, err := writer.Write([]byte("pong"))
    if err != nil {
        log.Errorf("ping failed: %s", err)
    }
}
