REPOSITORIES := icw_core_api icw_core_biz icw_activity_classification icw_activity_reasoning icw_activity_summary icw_common
IMAGE_REPOSITORIES := icw_core_api icw_core_biz icw_activity_classification icw_activity_reasoning icw_activity_summary icw_core_web

.DEFAULT_GOAL := help

.PHONY: help all verify image clean

help:
	@echo "Usage: make <COMMAND>"
	@echo ""
	@echo "Commands:"
	@echo "  make all     Run make all in all backend repositories"
	@echo "  make verify  Run make verify in all backend repositories"
	@echo "  make image   Run make image in all backend repositories"
	@echo "  make clean   Run make clean in all backend repositories"

all:
	@set -e; \
	for repo in $(REPOSITORIES); do \
		echo "cd $$repo && make all"; \
		$(MAKE) -C $$repo all; \
	done

verify:
	@set -e; \
	for repo in $(REPOSITORIES); do \
		echo "cd $$repo && make verify"; \
		$(MAKE) -C $$repo verify; \
	done

image:
	@set -e; \
	for repo in $(IMAGE_REPOSITORIES); do \
		echo "cd $$repo && make image"; \
		$(MAKE) -C $$repo image; \
	done

clean:
	@set -e; \
	for repo in $(REPOSITORIES); do \
		echo "cd $$repo && make clean"; \
		$(MAKE) -C $$repo clean; \
	done
