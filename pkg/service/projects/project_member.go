package projects

func (s *ProjectService) unAssignProjectMemberInJenkins(projectId, username string) error {
	for _, role := range AllRoleSlice {
		projectRole, err := s.Ds.Jenkins.GetProjectRole(GetProjectRoleName(projectId, role))
		if err != nil {
			return err
		}
		err = projectRole.UnAssignRole(username)
		if err != nil {
			return err
		}
		pipelineRole, err := s.Ds.Jenkins.GetProjectRole(GetPipelineRoleName(projectId, role))
		if err != nil {
			return err
		}
		err = pipelineRole.UnAssignRole(username)
		if err != nil {
			return err
		}
	}
	return nil
}
