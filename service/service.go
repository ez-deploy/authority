package service

import (
	"context"

	"github.com/ez-deploy/authority/db"
	pb "github.com/ez-deploy/protobuf/authority"
	"github.com/ez-deploy/protobuf/convert"
	"github.com/ez-deploy/protobuf/model"
	"github.com/wuhuizuo/sqlm"
)

type Service struct {
	AuthorityTable *sqlm.Table

	pb.UnimplementedOpsServer
}

func (s *Service) SetAuthorities(ctx context.Context, req *pb.Authorities) (*model.CommonResp, error) {
	insertItems := make([]interface{}, len(req.Authorities))

	for _, authority := range req.Authorities {
		rawResource, err := model.StringifyResource(authority.Resource)
		if err != nil {
			return model.NewCommonRespWithError(err), nil
		}

		insertItems = append(insertItems, &db.Authority{
			Identity: authority.Identity.Email,
			Resource: rawResource,
			Action:   authority.Action,
		})
	}

	if _, err := s.AuthorityTable.Inserts(insertItems); err != nil {
		return nil, err
	}

	return &model.CommonResp{}, nil
}

func (s *Service) ListAuthoritiesByIdentity(ctx context.Context, req *model.Identity) (*pb.ListAuthoritiesResp, error) {
	filter := sqlm.SelectorFilter{"identity": req.Email}
	listOptions := sqlm.ListOptions{AllColumns: true, OrderByColumn: "resource"}

	records, err := s.AuthorityTable.List(filter, listOptions)
	if err != nil {
		return nil, err
	}

	return newListAuthoritiesRespFromRecords(records)
}

func (s *Service) ListAuthoritiesByResource(ctx context.Context, req *model.Resource) (*pb.ListAuthoritiesResp, error) {
	rawResource, err := model.StringifyResource(req)
	if err != nil {
		return &pb.ListAuthoritiesResp{Error: model.NewError(err)}, nil
	}

	filter := sqlm.SelectorFilter{"resource": rawResource}
	listOptions := sqlm.ListOptions{AllColumns: true, OrderByColumn: "identity"}

	records, err := s.AuthorityTable.List(filter, listOptions)
	if err != nil {
		return nil, err
	}

	return newListAuthoritiesRespFromRecords(records)
}

func (s *Service) DeleteAuthorities(ctx context.Context, req *pb.Authorities) (*pb.DeleteAuthoritiesResp, error) {
	res := &pb.DeleteAuthoritiesResp{}

	for _, authority := range req.Authorities {
		rawResource, err := model.StringifyResource(authority.Resource)
		if err != nil {
			res.FailMessages = append(res.FailMessages, &pb.DeleteAuthoritiesResp_FailMessages{
				Error:     model.NewError(err),
				Authority: authority,
			})

			continue
		}

		filter := sqlm.SelectorFilter{
			"identity": authority.Identity.Email,
			"action":   authority.Action,
			"resource": rawResource,
		}
		if err := s.AuthorityTable.Delete(filter); err != nil {
			res.FailMessages = append(res.FailMessages, &pb.DeleteAuthoritiesResp_FailMessages{
				Error:     model.NewError(err),
				Authority: authority,
			})
		}
	}

	return res, nil
}

func newListAuthoritiesRespFromRecords(records []interface{}) (*pb.ListAuthoritiesResp, error) {
	res := &pb.ListAuthoritiesResp{Authorities: &pb.Authorities{}}

	for _, record := range records {
		dbAuthority := &db.Authority{}
		if err := convert.WithJSON(record, dbAuthority); err != nil {
			return nil, err
		}

		resource, err := model.NewResourceFromString(dbAuthority.Resource)
		if err != nil {
			return &pb.ListAuthoritiesResp{Error: model.NewError(err)}, nil
		}

		res.Authorities.Authorities = append(res.Authorities.Authorities, &model.Authority{
			Identity: &model.Identity{Email: dbAuthority.Identity},
			Action:   dbAuthority.Action,
			Resource: resource,
		})
	}

	return res, nil
}
