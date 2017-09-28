pkgname=zabbix-agent-extension-smartctl
pkgver=20170928.1_6e38cba
pkgrel=1
pkgdesc="Zabbix agent for SMART disks stats"
arch=('any')
license=('GPL')
depends=('smartmontools')
makedepends=('go')
install='install.sh'
source=("git+ssh://git@github.com/zarplata/$pkgname.git#branch=master")
md5sums=('SKIP')

pkgver() {
    cd "$srcdir/$pkgname"

    make ver
}
    
build() {
    cd "$srcdir/$pkgname"

    make
}

package() {
    cd "$srcdir/$pkgname"

    install -Dm 0755 .out/"${pkgname}" "${pkgdir}/usr/bin/${pkgname}"
    install -Dm 0644 "${pkgname}.conf" "${pkgdir}/etc/zabbix/zabbix_agentd.conf.d/${pkgname}.conf"
    install -Dm 0440 zabbix-smartctl "${pkgdir}/etc/sudoers.d/zabbix-smartctl"
}
