#version 330 core
out vec4 FragColor;

in vec3 Normal;  
in vec3 FragPos;
in vec3 vPos;
  
uniform vec3 lightPos; 
uniform vec3 lightColor;
uniform vec3 objectColor;

void main()
{
    // ambient
    float ambientStrength = 0.1;
    vec3 ambient = ambientStrength * lightColor;
  	
    // diffuse 
    vec3 norm = normalize(Normal);
    vec3 lightDir = normalize(lightPos - FragPos);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = diff * lightColor;

    vec3 result = vec3(0.0);
            
    if (vPos.y < 5) {   
        result = (ambient + diffuse) * vec3(0.19, 0.52, 0.25);
    }
    else if (vPos.y < 10) {   
        result = (ambient + diffuse) * vec3(0.45, 0.45, 0.45);
    }
    else { 
        result = (ambient + diffuse) * vec3(0.87, 0.87, 0.87);
    }

    FragColor = vec4(result, 1.0);
}
 
