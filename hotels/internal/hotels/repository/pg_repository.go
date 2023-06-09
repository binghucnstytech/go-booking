package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/hotels/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/pkg/hotels_errors"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/pkg/utils"
)

type hotelsPGRepository struct {
	db *pgxpool.Pool
}

func NewHotelsPGRepository(db *pgxpool.Pool) *hotelsPGRepository {
	return &hotelsPGRepository{db: db}
}

func (h *hotelsPGRepository) CreateHotel(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "hotelsPGRepository.CreateHotel")
	defer span.Finish()

	point := utils.GeneratePointToGeoFromFloat64(*hotel.Latitude, *hotel.Longitude)

	var res models.Hotel
	if err := h.db.QueryRow(
		ctx,
		createHotelQuery,
		hotel.Name,
		hotel.Location,
		hotel.Description,
		hotel.Image,
		hotel.Photos,
		point,
		hotel.Email,
		hotel.Country,
		hotel.City,
		hotel.Rating,
	).Scan(&res.HotelID, &res.CreatedAt, &res.UpdatedAt); err != nil {
		return nil, errors.Wrap(err, "CreateHotel.Scan")
	}

	hotel.HotelID = res.HotelID
	hotel.CreatedAt = res.CreatedAt
	hotel.UpdatedAt = res.UpdatedAt

	return hotel, nil
}

// UpdateHotel update single hotel
func (h *hotelsPGRepository) UpdateHotel(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "hotelsPGRepository.UpdateHotel")
	defer span.Finish()

	point := utils.GeneratePointToGeoFromFloat64(*hotel.Latitude, *hotel.Longitude)
	var res models.Hotel
	if err := h.db.QueryRow(
		ctx,
		updateHotelQuery,
		hotel.Email,
		hotel.Name,
		hotel.Location,
		hotel.Description,
		hotel.Country,
		hotel.City,
		point,
		hotel.HotelID,
	).Scan(
		&res.HotelID,
		&res.Email,
		&res.Name,
		&res.Location,
		&res.Description,
		&res.CommentsCount,
		&res.Country,
		&res.City,
		&res.Latitude,
		&res.Longitude,
		&res.Rating,
		&res.Photos,
		&res.Image,
		&res.CreatedAt,
		&res.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "db.QueryRow.Scan")
	}

	return &res, nil
}

// GetHotelByID get hotel by id
func (h *hotelsPGRepository) GetHotelByID(ctx context.Context, hotelID uuid.UUID) (*models.Hotel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "hotelsPGRepository.GetHotelByID")
	defer span.Finish()

	var hotel models.Hotel
	if err := h.db.QueryRow(ctx, getHotelByIDQuery, hotelID).Scan(
		&hotel.HotelID,
		&hotel.Email,
		&hotel.Name,
		&hotel.Location,
		&hotel.Description,
		&hotel.CommentsCount,
		&hotel.Country,
		&hotel.City,
		&hotel.Latitude,
		&hotel.Longitude,
		&hotel.Rating,
		&hotel.Photos,
		&hotel.Image,
		&hotel.CreatedAt,
		&hotel.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "db.QueryRow.Scan")
	}

	return &hotel, nil
}

// GetHotels
func (h *hotelsPGRepository) GetHotels(ctx context.Context, query *utils.PaginationQuery) (*models.HotelsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "hotelsPGRepository.GetHotels")
	defer span.Finish()

	var total int
	if err := h.db.QueryRow(ctx, getTotalHotelsCountQuery).Scan(&total); err != nil {
		return nil, errors.Wrap(err, "db.Query")
	}
	if total == 0 {
		return &models.HotelsList{
			TotalCount: total,
			TotalPages: 0,
			Page:       0,
			Size:       0,
			HasMore:    false,
			Hotels:     make([]*models.Hotel, 0),
		}, nil
	}

	rows, err := h.db.Query(ctx, getHotelsQuery, query.GetOffset(), query.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "db.Query")
	}
	defer rows.Close()

	hotels := make([]*models.Hotel, 0, query.GetLimit())
	for rows.Next() {
		var hotel models.Hotel
		if err := rows.Scan(
			&hotel.HotelID,
			&hotel.Email,
			&hotel.Name,
			&hotel.Location,
			&hotel.Description,
			&hotel.CommentsCount,
			&hotel.Country,
			&hotel.City,
			&hotel.Latitude,
			&hotel.Longitude,
			&hotel.Rating,
			&hotel.Photos,
			&hotel.Image,
			&hotel.CreatedAt,
			&hotel.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}
		hotels = append(hotels, &hotel)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err")
	}

	log.Printf("HOTELS: %-v", hotels)

	return &models.HotelsList{
		TotalCount: total,
		TotalPages: query.GetTotalPages(total),
		Page:       query.GetPage(),
		Size:       query.GetSize(),
		HasMore:    query.GetHasMore(total),
		Hotels:     hotels,
	}, nil
}

// UpdateHotelImage
func (h *hotelsPGRepository) UpdateHotelImage(ctx context.Context, hotelID uuid.UUID, imageURL string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "hotelsPGRepository.UpdateHotelImage")
	defer span.Finish()

	updateHotelImageQuery := `UPDATE hotels SET image = $1 WHERE hotel_id = $2`

	result, err := h.db.Exec(ctx, updateHotelImageQuery, imageURL, hotelID)
	if err != nil {
		return errors.Wrap(err, "db.Exec")
	}

	if result.RowsAffected() == 0 {
		return errors.Wrap(hotels_errors.ErrHotelNotFound, "RowsAffected")
	}

	return nil
}
