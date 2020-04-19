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

package doc

import (
	"net/http"

	"github.com/rebel-l/smis"
)

// Init initialises the doc endpoints
func Init(svc *smis.Service) error {
	_, err := svc.RegisterFileServer("/doc", http.MethodGet, "endpoint/doc/web")
	return err
}
