Pod::Spec.new do |spec|
  spec.name         = 'Gnod'
  spec.version      = '{{.Version}}'
  spec.license      = { :type => 'GNU Lesser General Public License, Version 3.0' }
  spec.homepage     = 'https://github.com/neworder-network/go-neworder'
  spec.authors      = { {{range .Contributors}}
		'{{.Name}}' => '{{.Email}}',{{end}}
	}
  spec.summary      = 'iOS Neworder Client'
  spec.source       = { :git => 'https://github.com/neworder-network/go-neworder.git', :commit => '{{.Commit}}' }

	spec.platform = :ios
  spec.ios.deployment_target  = '9.0'
	spec.ios.vendored_frameworks = 'Frameworks/Gnod.framework'

	spec.prepare_command = <<-CMD
    curl https://gnodstore.blob.core.windows.net/builds/{{.Archive}}.tar.gz | tar -xvz
    mkdir Frameworks
    mv {{.Archive}}/Gnod.framework Frameworks
    rm -rf {{.Archive}}
  CMD
end
