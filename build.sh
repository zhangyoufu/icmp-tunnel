#!/bin/sh

cd -- "$(dirname -- "$0")"

export GOOS=linux
if [ "$TARGETOS" != "$GOOS" ]; then
	echo 'This program only supports Linux.'
	exit 1
fi

# https://github.com/golang/go/blob/0f6ee42fe063a48d7825bc03097bbb714aafdb7d/test/run.go#L1599-L1613
export GOARCH=${TARGETARCH}
case $TARGETARCH in
	386)
		# defaults to sse2 if unset
		export GO386=$TARGETVARIANT
	;;
	amd64)
		# defaults to v1 if unset
		export GOAMD64=$TARGETVARIANT
	;;
	arm)
		# https://github.com/containerd/containerd/blob/4902059cb554f4f06a8d06a12134c17117809f4e/platforms/cpuinfo.go#L113-L128
		# https://github.com/golang/go/wiki/GoArm
		case $TARGETVARIANT in
			'')
				# default value determined by xgetgoarm()
			;;
			v7|v6|v5)
				export GOARM=${TARGETVARIANT#v}
			;;
			v4|v3)
				echo "unsupported TARGETVARIANT=$TARGETVARIANT for TARGETARCH=$TARGETARCH"
				exit 1
			;;
			*)
				echo "unknown TARGETVARIANT=$TARGETVARIANT for TARGETARCH=$TARGETARCH"
				exit 1
			;;
		esac
	;;
	mips|mipsle)
		# defaults to hardfloat if unset
		export GOMIPS=$TARGETVARIANT
	;;
	mips64|mips64le)
		# defaults to hardfloat if unset
		export GOMIPS64=$TARGETVARIANT
	;;
	ppc64|ppc64le)
		# defaults to power8 if unset
		export GOPPC64=$TARGETVARIANT
	;;
	arm64|s390x|riscv64|loong64)
		if [ -n "$TARGETVARIANT" ]; then
			echo "unknown TARGETVARIANT=$TARGETVARIANT for TARGETARCH=$TARGETARCH"
			exit 1
		fi
	;;
	*)
esac

export CGO_ENABLED=0

exec go build -ldflags="-s -w" -trimpath -o /icmp-tunnel .
