all: clone update cleanup

clone:
	git clone -b DEV-15047-slim-cnvrg  git@github.com:AccessibleAI/cnvrg-operator.git

rename:
	sed -i .bak 's/name: cnvrg/name: cnvrg-operator/g' cnvrg-operator/charts/operator/Chart.yaml

update:
	helm dependency update ./slim-operator-delivery

cleanup:
	rm -rf cnvrg-operator
