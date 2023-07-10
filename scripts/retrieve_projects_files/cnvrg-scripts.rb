def get_file_names_for_all_projects()
	i = 0;
	projects  = Project.where(is_git: false)
	p_to_files = {}

	for p in projects
		i += 1
		slug_ = p.slug
		if(p_to_files.key?(p.slug))
			slug_ += i.to_s
		end;
		p_to_files[slug_] = get_file_names(get_bvs(p))
		puts i, p.slug + "   ", p_to_files.key?(p.slug)
	end;
	return p_to_files
end;

def get_bvs(proj)
	lc = proj.commit_ids.last
	ctbv = CommitToBlobVersion.where(commit_id: lc)
	bv_ids = ctbv.map { |temp| temp.blob_v_id}
	bvs = BlobVersion.where(id: bv_ids)
	return bvs
end;

def get_file_names(bvs)
	file_names = []
	for bv in bvs
		file_names.append(bv.external_attachment_path)
	end;
	return file_names
end;


i = 0
project_files = {}
project_files.each do |p_slug, files|
    puts i
	i+=1
  	File.write("projects_to_files.txt", p_slug + " = [", mode: "a")
  	for f_ in files
		if(f_)
			File.write("projects_to_files.txt", '"' + f_ + '"' + ',', mode: "a")
		end;
	end;
    File.write("projects_to_files.txt", ']\n', mode: "a")
end;

