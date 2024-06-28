package metadata

import (
	"context"

	"google.golang.org/grpc/metadata"
)

type Payload struct {
	UserID   string
	UserName string
	Ip       string
	Token    string
	Role     string
}

const (
	MDUserIDKey     = "user_id"
	MDUserNameKey   = "user_name"
	MDIpKey         = "ip"
	MDTokenKey      = "token"
	MDRoleKey       = "role"
	MDXForwardedFor = "x-forwarded-for"
)

// ImportUserInfoToCtx creates metadata containing user information from the provided payload.
//
// It takes a payload struct as input and returns metadata.MD.
func ImportUserInfoToCtx(payload *Payload) metadata.MD {
	md := metadata.Pairs(MDUserIDKey, payload.UserID)
	md.Append(MDUserNameKey, payload.UserName)
	md.Append(MDTokenKey, payload.Token)
	md.Append(MDRoleKey, payload.Role)
	md.Append(MDIpKey, payload.Ip)

	return md
}

// ExtractUserInfoFromCtx extracts user information from the context.
func ExtractUserInfoFromCtx(ctx context.Context) (*Payload, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, false
	}

	vals := []string{}

	for _, key := range []string{MDUserIDKey, MDUserNameKey, MDIpKey, MDTokenKey, MDRoleKey, MDXForwardedFor} {
		values := md.Get(key)
		if len(values) < 1 {
			vals = append(vals, "")
		} else {
			vals = append(vals, values[0])
		}
	}

	ip := vals[2]
	if vals[2] == "" {
		ip = vals[4]
	}

	return &Payload{
		UserID:   vals[0],
		UserName: vals[1],
		Ip:       ip,
		Token:    vals[3],
		Role:     vals[4],
	}, true
}

// InjectIncomingCtxToOutgoingCtx injects metadata from incoming context to outgoing context.
//
// It takes a context.Context parameter and returns a context.Context.
func InjectIncomingCtxToOutgoingCtx(ctx context.Context) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)

	return metadata.NewOutgoingContext(ctx, md)
}
