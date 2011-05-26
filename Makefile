include $(GOROOT)/src/Make.inc

TARG=gourl

O_FILES = gourl.6

all:
	make clean
	make $(TARG)

$(TARG): $(O_FILES)
	$(LD) -o $@ $(O_FILES)
	@echo "Done. Executable is: $@"

$(O_FILES): %.6: %.go
	$(GC) -c $<

clean:
	rm -rf *.[$(OS)o] *.a [$(OS)].out _obj $(TARG) *.6
