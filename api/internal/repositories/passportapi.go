package repositories

import (
	"em-test/internal/config"
	"em-test/internal/lib/dto"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// var _ services.UserFinder = (*PassportApi)(nil)

type PassportApi struct {
	host string
}

func NewPassportApi(config *config.Config) *PassportApi {
	return &PassportApi{
		host: config.PassportApi.Host,
	}
}

func (p *PassportApi) GetInfo(serie int, number int) (*dto.UserInfoDto, error) {

	fn := "PassportApi.GetInfo"
	logger := slog.With(slog.String("fn", fn))
	endpoint := fmt.Sprintf("%s/info?passportSerie=%d&passportNumber=%d", p.host, serie, number)

	logger.Debug("sending request", slog.String("endpoint", endpoint))

	response, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	info := &dto.UserInfoDto{}

	if err := json.NewDecoder(response.Body).Decode(info); err != nil {
		logger.Error("cannot parse response body", slog.String("err", err.Error()))
		return nil, err
	}

	return info, nil
}
